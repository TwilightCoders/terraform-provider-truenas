// Package middleware is a JSON-RPC 2.0 client for the TrueNAS versioned websocket API.
//
// The client pins one API version (e.g. /api/v25.10.5), authenticates with an API key via
// auth.login_ex, and enables new-style jobs so that job-backed methods return their final
// result in the ordinary response. It reconnects lazily: a dropped connection fails the calls
// that were in flight and the next call dials again.
package middleware

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"net/http"
	"sync"
	"sync/atomic"
	"time"

	"github.com/coder/websocket"
)

const (
	defaultMaxConcurrency = 8
	defaultReadLimit      = 64 << 20
	defaultKeepAlive      = 30 * time.Second
	defaultWriteTimeout   = 30 * time.Second
	defaultBusyRetries    = 5
	defaultBusyBackoff    = 200 * time.Millisecond
)

// Config configures a Client.
type Config struct {
	// Host is the TrueNAS web UI address, optionally with a port: "nas.lan" or "nas.lan:8443".
	Host string
	// APIVersion is the newest API version the caller understands, e.g. "v25.10.5". Dial picks the
	// highest version the server offers in the same major.minor series that is not newer.
	APIVersion string
	// Username owns the API key.
	Username string
	// APIKey authenticates the session.
	APIKey string
	// TLS configures certificate verification. Ignored when HTTPClient is set.
	TLS TLSConfig
	// HTTPClient overrides the client used for the handshake and version discovery.
	HTTPClient *http.Client
	// MaxConcurrency bounds in-flight calls. Default 8.
	MaxConcurrency int
	// ReadLimit bounds a single response in bytes. Default 64 MiB.
	ReadLimit int64
	// KeepAlive is the ping interval. Default 30s; negative disables pings.
	KeepAlive time.Duration
	// Logger receives call-level debug logs. Params are never logged.
	Logger Logger
}

// Logger receives structured debug logs.
type Logger interface {
	Debug(ctx context.Context, msg string, fields map[string]any)
}

// Client is a TrueNAS middleware client. It is safe for concurrent use.
type Client struct {
	cfg      Config
	httpc    *http.Client
	version  string
	endpoint string
	sem      chan struct{}
	nextID   atomic.Uint64

	mu     sync.Mutex
	sess   *session
	closed bool
}

// Dial negotiates the API version, connects and authenticates.
func Dial(ctx context.Context, cfg Config) (*Client, error) {
	if err := cfg.validate(); err != nil {
		return nil, err
	}
	cfg.applyDefaults()

	httpc := cfg.HTTPClient
	if httpc == nil {
		tlsCfg, err := cfg.TLS.build()
		if err != nil {
			return nil, err
		}
		httpc = &http.Client{Transport: &http.Transport{TLSClientConfig: tlsCfg, Proxy: http.ProxyFromEnvironment}}
	}

	available, err := fetchVersions(ctx, httpc, cfg.Host)
	if err != nil {
		return nil, err
	}
	version, err := SelectVersion(available, cfg.APIVersion)
	if err != nil {
		return nil, err
	}

	c := &Client{
		cfg:      cfg,
		httpc:    httpc,
		version:  version,
		endpoint: fmt.Sprintf("wss://%s/api/%s", cfg.Host, version),
		sem:      make(chan struct{}, cfg.MaxConcurrency),
	}
	if _, err := c.session(ctx); err != nil {
		return nil, err
	}
	return c, nil
}

// APIVersion returns the negotiated API version.
func (c *Client) APIVersion() string { return c.version }

// Call invokes method with positional params and returns the raw JSON result.
func (c *Client) Call(ctx context.Context, method string, params ...any) (json.RawMessage, error) {
	select {
	case c.sem <- struct{}{}:
		defer func() { <-c.sem }()
	case <-ctx.Done():
		return nil, ctx.Err()
	}

	start := time.Now()
	result, err := c.callWithBusyRetry(ctx, method, params)
	c.log(ctx, "truenas call", map[string]any{
		"method":      method,
		"duration_ms": time.Since(start).Milliseconds(),
		"error":       errString(err),
	})
	return result, err
}

// CallInto invokes method and decodes the result into out.
func (c *Client) CallInto(ctx context.Context, method string, out any, params ...any) error {
	raw, err := c.Call(ctx, method, params...)
	if err != nil {
		return err
	}
	if err := json.Unmarshal(raw, out); err != nil {
		return fmt.Errorf("%s: decoding result: %w", method, err)
	}
	return nil
}

// Close terminates the connection. Calls after Close fail.
func (c *Client) Close() error {
	c.mu.Lock()
	defer c.mu.Unlock()
	c.closed = true
	if c.sess == nil {
		return nil
	}
	err := c.sess.close()
	c.sess = nil
	return err
}

func (c *Client) callWithBusyRetry(ctx context.Context, method string, params []any) (json.RawMessage, error) {
	backoff := defaultBusyBackoff
	for attempt := 0; ; attempt++ {
		s, err := c.session(ctx)
		if err != nil {
			var authErr *AuthError
			if errors.As(err, &authErr) {
				return nil, err
			}
			return nil, &ConnectionError{Method: method, Err: err}
		}
		result, err := s.call(ctx, c.nextID.Add(1), method, params)
		var e *Error
		if attempt < defaultBusyRetries && errors.As(err, &e) && e.Code == CodeTooManyConcurrentCalls {
			if werr := sleep(ctx, backoff); werr != nil {
				return nil, werr
			}
			backoff *= 2
			continue
		}
		return result, err
	}
}

// session returns the live session, dialing a new one when needed.
func (c *Client) session(ctx context.Context) (*session, error) {
	c.mu.Lock()
	defer c.mu.Unlock()
	if c.closed {
		return nil, errors.New("truenas client is closed")
	}
	if c.sess != nil && c.sess.alive() {
		return c.sess, nil
	}
	s, err := c.connect(ctx)
	if err != nil {
		return nil, err
	}
	c.sess = s
	return s, nil
}

func (c *Client) connect(ctx context.Context) (*session, error) {
	conn, resp, err := websocket.Dial(ctx, c.endpoint, &websocket.DialOptions{HTTPClient: c.httpc})
	if resp != nil && resp.Body != nil {
		_ = resp.Body.Close()
	}
	if err != nil {
		return nil, fmt.Errorf("connecting to %s: %w", c.endpoint, err)
	}
	conn.SetReadLimit(c.cfg.ReadLimit)
	s := newSession(conn, c.cfg.KeepAlive)

	// New-style jobs: job methods reply once the job finishes, with its result or error.
	if _, err := s.call(ctx, c.nextID.Add(1), "core.set_options", []any{map[string]any{"legacy_jobs": false}}); err != nil {
		_ = s.close()
		return nil, err
	}

	var login struct {
		ResponseType string `json:"response_type"`
	}
	raw, err := s.call(ctx, c.nextID.Add(1), "auth.login_ex", []any{map[string]any{
		"mechanism": "API_KEY_PLAIN",
		"username":  c.cfg.Username,
		"api_key":   c.cfg.APIKey,
	}})
	if err == nil {
		err = json.Unmarshal(raw, &login)
	}
	if err != nil {
		_ = s.close()
		return nil, fmt.Errorf("authenticating: %w", err)
	}
	if login.ResponseType != "SUCCESS" {
		_ = s.close()
		return nil, &AuthError{ResponseType: login.ResponseType}
	}
	return s, nil
}

func (c *Client) log(ctx context.Context, msg string, fields map[string]any) {
	if c.cfg.Logger != nil {
		c.cfg.Logger.Debug(ctx, msg, fields)
	}
}

func (cfg *Config) validate() error {
	var missing []string
	if cfg.Host == "" {
		missing = append(missing, "host")
	}
	if cfg.APIVersion == "" {
		missing = append(missing, "API version")
	}
	if cfg.Username == "" {
		missing = append(missing, "username")
	}
	if cfg.APIKey == "" {
		missing = append(missing, "API key")
	}
	if len(missing) > 0 {
		return fmt.Errorf("truenas client config is missing: %v", missing)
	}
	return nil
}

func (cfg *Config) applyDefaults() {
	if cfg.MaxConcurrency <= 0 {
		cfg.MaxConcurrency = defaultMaxConcurrency
	}
	if cfg.ReadLimit <= 0 {
		cfg.ReadLimit = defaultReadLimit
	}
	if cfg.KeepAlive == 0 {
		cfg.KeepAlive = defaultKeepAlive
	}
}

func sleep(ctx context.Context, d time.Duration) error {
	t := time.NewTimer(d)
	defer t.Stop()
	select {
	case <-t.C:
		return nil
	case <-ctx.Done():
		return ctx.Err()
	}
}

func errString(err error) string {
	if err == nil {
		return ""
	}
	return err.Error()
}
