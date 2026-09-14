package middleware_test

import (
	"context"
	"crypto/sha256"
	"encoding/hex"
	"encoding/json"
	"encoding/pem"
	"errors"
	"net/http"
	"net/http/httptest"
	"strings"
	"sync"
	"sync/atomic"
	"testing"
	"time"

	"github.com/TwilightCoders/terraform-provider-truenas/internal/middleware"
	"github.com/TwilightCoders/terraform-provider-truenas/internal/middleware/middlewaretest"
)

func dial(t *testing.T, srv *middlewaretest.Server, mutate ...func(*middleware.Config)) *middleware.Client {
	t.Helper()
	cfg := srv.Config()
	for _, m := range mutate {
		m(&cfg)
	}
	c, err := middleware.Dial(context.Background(), cfg)
	if err != nil {
		t.Fatalf("Dial: %v", err)
	}
	t.Cleanup(func() { _ = c.Close() })
	return c
}

func TestCallReturnsResultAndSendsPositionalParams(t *testing.T) {
	srv := middlewaretest.NewServer(t)
	srv.Handle("cronjob.update", func(_ context.Context, params []json.RawMessage) (any, error) {
		return map[string]any{"id": 7, "params": len(params)}, nil
	})
	c := dial(t, srv)

	var out struct {
		ID     int `json:"id"`
		Params int `json:"params"`
	}
	if err := c.CallInto(context.Background(), "cronjob.update", &out, 7, map[string]any{"enabled": false}); err != nil {
		t.Fatalf("CallInto: %v", err)
	}
	if out.ID != 7 || out.Params != 2 {
		t.Fatalf("unexpected result %+v", out)
	}

	calls := srv.Calls()
	if len(calls) != 1 || calls[0].Method != "cronjob.update" {
		t.Fatalf("calls = %+v", calls)
	}
	if got := string(calls[0].Params[0]); got != "7" {
		t.Errorf("first param = %s, want 7", got)
	}
	if got := string(calls[0].Params[1]); got != `{"enabled":false}` {
		t.Errorf("second param = %s", got)
	}
	if c.APIVersion() != middlewaretest.APIVersion {
		t.Errorf("APIVersion = %q", c.APIVersion())
	}
}

func TestCallWithoutParamsSendsEmptyArray(t *testing.T) {
	srv := middlewaretest.NewServer(t)
	srv.Handle("smb.config", func(_ context.Context, params []json.RawMessage) (any, error) {
		if params == nil {
			return nil, errors.New("params missing")
		}
		return map[string]any{"id": 1}, nil
	})
	c := dial(t, srv)
	if _, err := c.Call(context.Background(), "smb.config"); err != nil {
		t.Fatal(err)
	}
}

func TestCallIntoDecodeError(t *testing.T) {
	srv := middlewaretest.NewServer(t)
	srv.Handle("system.version", func(context.Context, []json.RawMessage) (any, error) { return "25.10.7", nil })
	c := dial(t, srv)
	var out int
	err := c.CallInto(context.Background(), "system.version", &out)
	if err == nil || !strings.Contains(err.Error(), "decoding result") {
		t.Fatalf("err = %v", err)
	}
}

func TestErrorMapping(t *testing.T) {
	srv := middlewaretest.NewServer(t)
	srv.Handle("cronjob.get_instance", func(context.Context, []json.RawMessage) (any, error) {
		return nil, &middleware.Error{Errno: 2, Errname: "ENOENT", Reason: "cronjob 9 does not exist"}
	})
	srv.Handle("sharing.smb.create", func(context.Context, []json.RawMessage) (any, error) {
		return nil, &middleware.Error{Errno: 22, Errname: "EINVAL", Fields: []middleware.FieldError{
			{Attribute: "sharing_smb_create.path", Message: "Path must exist", Errno: 2},
			{Attribute: "sharing_smb_create.name", Message: "Share name conflicts", Errno: 17},
		}}
	})
	srv.Handle("app.create", func(context.Context, []json.RawMessage) (any, error) {
		return nil, errors.New("boom")
	})
	c := dial(t, srv)
	ctx := context.Background()

	_, err := c.Call(ctx, "cronjob.get_instance", 9)
	if !middleware.IsNotFound(err) {
		t.Fatalf("IsNotFound(%v) = false", err)
	}
	if want := "cronjob.get_instance: cronjob 9 does not exist [ENOENT]"; err.Error() != want {
		t.Errorf("Error() = %q, want %q", err.Error(), want)
	}
	if middleware.IsValidation(err) || middleware.IsRetryable(err) {
		t.Error("not-found error misclassified")
	}

	_, err = c.Call(ctx, "sharing.smb.create", map[string]any{})
	var e *middleware.Error
	if !errors.As(err, &e) || !middleware.IsValidation(err) {
		t.Fatalf("expected validation error, got %v", err)
	}
	if e.Code != middleware.CodeInvalidParams || len(e.Fields) != 2 || e.Fields[1].Attribute != "sharing_smb_create.name" || e.Fields[1].Errno != 17 {
		t.Errorf("fields = %+v code=%d", e.Fields, e.Code)
	}

	_, err = c.Call(ctx, "app.create", map[string]any{})
	if !errors.As(err, &e) || e.Errname != "EINVAL" || e.Reason != "boom" || e.Code != middleware.CodeCallError {
		t.Errorf("generic error = %#v", err)
	}

	_, err = c.Call(ctx, "no.such.method")
	if !errors.As(err, &e) || e.Code != middleware.CodeMethodNotFound {
		t.Fatalf("unknown method error = %v", err)
	}
	if want := "no.such.method: Method does not exist"; err.Error() != want {
		t.Errorf("Error() = %q, want %q", err.Error(), want)
	}
}

func TestDialRejectsBadCredentials(t *testing.T) {
	srv := middlewaretest.NewServer(t)
	cfg := srv.Config()
	cfg.APIKey = "wrong"
	_, err := middleware.Dial(context.Background(), cfg)
	var authErr *middleware.AuthError
	if !errors.As(err, &authErr) || authErr.ResponseType != "AUTH_ERR" {
		t.Fatalf("err = %v", err)
	}
	if !strings.Contains(err.Error(), "invalid username or API key") {
		t.Errorf("message = %q", err.Error())
	}
}

func TestAuthErrorMessages(t *testing.T) {
	for rt, want := range map[string]string{
		"EXPIRED":      "expired",
		"OTP_REQUIRED": "one-time password",
		"REDIRECT":     "REDIRECT",
	} {
		srv := middlewaretest.NewServer(t)
		srv.SetLoginResponse(rt)
		_, err := middleware.Dial(context.Background(), srv.Config())
		if err == nil || !strings.Contains(err.Error(), want) {
			t.Errorf("%s: err = %v", rt, err)
		}
	}
}

func TestDialValidatesConfig(t *testing.T) {
	_, err := middleware.Dial(context.Background(), middleware.Config{})
	if err == nil {
		t.Fatal("expected error")
	}
	for _, field := range []string{"host", "API version", "username", "API key"} {
		if !strings.Contains(err.Error(), field) {
			t.Errorf("error %q does not mention %s", err, field)
		}
	}
}

func TestDialNegotiatesVersion(t *testing.T) {
	srv := middlewaretest.NewServer(t)
	srv.SetVersions("v25.04.2", "v25.10.2", "v25.10.3", "v26.0.0")
	c := dial(t, srv)
	if c.APIVersion() != "v25.10.3" {
		t.Fatalf("APIVersion = %q, want v25.10.3", c.APIVersion())
	}

	srv.SetVersions("v25.04.2")
	_, err := middleware.Dial(context.Background(), srv.Config())
	if err == nil || !strings.Contains(err.Error(), "needs v25.10.x") {
		t.Fatalf("err = %v", err)
	}
}

func TestDialVersionDiscoveryFailures(t *testing.T) {
	notFound := httptest.NewTLSServer(http.NotFoundHandler())
	defer notFound.Close()
	garbage := httptest.NewTLSServer(http.HandlerFunc(func(w http.ResponseWriter, _ *http.Request) { _, _ = w.Write([]byte("<html>")) }))
	defer garbage.Close()

	for name, srv := range map[string]*httptest.Server{"404": notFound, "garbage": garbage} {
		cfg := middleware.Config{Host: srv.Listener.Addr().String(), APIVersion: "v25.10.5", Username: "u", APIKey: "k", HTTPClient: srv.Client()}
		if _, err := middleware.Dial(context.Background(), cfg); err == nil || !strings.Contains(err.Error(), "discovering API versions") {
			t.Errorf("%s: err = %v", name, err)
		}
	}

	cfg := middleware.Config{Host: "127.0.0.1:1", APIVersion: "v25.10.5", Username: "u", APIKey: "k"}
	if _, err := middleware.Dial(context.Background(), cfg); err == nil {
		t.Error("expected connection refused")
	}
}

func TestDialFailsWhenWebsocketUnavailable(t *testing.T) {
	versionsOnly := httptest.NewTLSServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.URL.Path == "/api/versions" {
			_, _ = w.Write([]byte(`["v25.10.5"]`))
			return
		}
		http.NotFound(w, r)
	}))
	defer versionsOnly.Close()
	cfg := middleware.Config{Host: versionsOnly.Listener.Addr().String(), APIVersion: "v25.10.5", Username: "u", APIKey: "k", HTTPClient: versionsOnly.Client()}
	if _, err := middleware.Dial(context.Background(), cfg); err == nil || !strings.Contains(err.Error(), "connecting to") {
		t.Fatalf("err = %v", err)
	}
}

func TestBusyServerIsRetried(t *testing.T) {
	srv := middlewaretest.NewServer(t)
	var attempts atomic.Int32
	srv.Handle("pool.query", func(context.Context, []json.RawMessage) (any, error) {
		if attempts.Add(1) < 3 {
			return nil, &middleware.Error{Code: middleware.CodeTooManyConcurrentCalls, Message: "Maximum number of concurrent calls (20) has exceeded"}
		}
		return []any{}, nil
	})
	c := dial(t, srv)
	if _, err := c.Call(context.Background(), "pool.query"); err != nil {
		t.Fatalf("Call: %v", err)
	}
	if attempts.Load() != 3 {
		t.Errorf("attempts = %d, want 3", attempts.Load())
	}
}

func TestBusyRetryStopsOnContextCancel(t *testing.T) {
	srv := middlewaretest.NewServer(t)
	srv.Handle("pool.query", func(context.Context, []json.RawMessage) (any, error) {
		return nil, &middleware.Error{Code: middleware.CodeTooManyConcurrentCalls, Message: "busy"}
	})
	c := dial(t, srv)
	ctx, cancel := context.WithTimeout(context.Background(), 50*time.Millisecond)
	defer cancel()
	if _, err := c.Call(ctx, "pool.query"); !errors.Is(err, context.DeadlineExceeded) {
		t.Fatalf("err = %v", err)
	}
}

func TestConnectionDropFailsInFlightCallAndReconnects(t *testing.T) {
	srv := middlewaretest.NewServer(t)
	started := make(chan struct{})
	srv.Handle("cronjob.create", func(ctx context.Context, _ []json.RawMessage) (any, error) {
		close(started)
		<-ctx.Done()
		return nil, ctx.Err()
	})
	srv.Handle("cronjob.query", func(context.Context, []json.RawMessage) (any, error) { return []any{}, nil })
	c := dial(t, srv)

	errc := make(chan error, 1)
	go func() {
		_, err := c.Call(context.Background(), "cronjob.create", map[string]any{})
		errc <- err
	}()
	<-started
	srv.DropConnections()

	err := <-errc
	var ce *middleware.ConnectionError
	if !errors.As(err, &ce) || !ce.Sent {
		t.Fatalf("err = %v, want ConnectionError with Sent", err)
	}
	if middleware.IsRetryable(err) {
		t.Error("a sent call must not be retryable")
	}
	if !strings.Contains(err.Error(), "outcome unknown") {
		t.Errorf("message = %q", err.Error())
	}

	if _, err := c.Call(context.Background(), "cronjob.query"); err != nil {
		t.Fatalf("call after reconnect: %v", err)
	}
}

func TestReconnectFailureIsRetryable(t *testing.T) {
	srv := middlewaretest.NewServer(t)
	srv.Handle("cronjob.query", func(context.Context, []json.RawMessage) (any, error) { return []any{}, nil })
	c := dial(t, srv)

	srv.SetVersions() // the websocket route now 404s
	srv.DropConnections()
	time.Sleep(50 * time.Millisecond)

	_, err := c.Call(context.Background(), "cronjob.query")
	var ce *middleware.ConnectionError
	if !errors.As(err, &ce) || ce.Sent || !middleware.IsRetryable(err) {
		t.Fatalf("err = %v, want unsent ConnectionError", err)
	}
	if !strings.Contains(err.Error(), "before the request was sent") {
		t.Errorf("message = %q", err.Error())
	}
}

func TestReconnectWithRevokedKeyReturnsAuthError(t *testing.T) {
	srv := middlewaretest.NewServer(t)
	c := dial(t, srv)
	srv.SetLoginResponse("EXPIRED")
	srv.DropConnections()
	time.Sleep(50 * time.Millisecond)

	_, err := c.Call(context.Background(), "cronjob.query")
	var authErr *middleware.AuthError
	if !errors.As(err, &authErr) {
		t.Fatalf("err = %v", err)
	}
}

func TestCallContextCancelLeavesConnectionUsable(t *testing.T) {
	srv := middlewaretest.NewServer(t)
	release := make(chan struct{})
	srv.Handle("slow", func(context.Context, []json.RawMessage) (any, error) {
		<-release
		return 1, nil
	})
	srv.Handle("fast", func(context.Context, []json.RawMessage) (any, error) { return 2, nil })
	c := dial(t, srv)
	defer close(release)

	ctx, cancel := context.WithTimeout(context.Background(), 30*time.Millisecond)
	defer cancel()
	if _, err := c.Call(ctx, "slow"); !errors.Is(err, context.DeadlineExceeded) {
		t.Fatalf("err = %v", err)
	}
	raw, err := c.Call(context.Background(), "fast")
	if err != nil || string(raw) != "2" {
		t.Fatalf("fast = %s, %v", raw, err)
	}
}

func TestConcurrencyIsBounded(t *testing.T) {
	srv := middlewaretest.NewServer(t)
	var inFlight, peak atomic.Int32
	srv.Handle("pool.dataset.query", func(context.Context, []json.RawMessage) (any, error) {
		n := inFlight.Add(1)
		defer inFlight.Add(-1)
		for {
			p := peak.Load()
			if n <= p || peak.CompareAndSwap(p, n) {
				break
			}
		}
		time.Sleep(5 * time.Millisecond)
		return []any{}, nil
	})
	c := dial(t, srv, func(cfg *middleware.Config) { cfg.MaxConcurrency = 3 })

	var wg sync.WaitGroup
	for range 30 {
		wg.Go(func() {
			if _, err := c.Call(context.Background(), "pool.dataset.query"); err != nil {
				t.Error(err)
			}
		})
	}
	wg.Wait()
	if peak.Load() > 3 {
		t.Errorf("peak concurrency = %d, want <= 3", peak.Load())
	}
}

func TestCallWaitsForSlotRespectingContext(t *testing.T) {
	srv := middlewaretest.NewServer(t)
	release := make(chan struct{})
	srv.Handle("slow", func(context.Context, []json.RawMessage) (any, error) {
		<-release
		return 1, nil
	})
	c := dial(t, srv, func(cfg *middleware.Config) { cfg.MaxConcurrency = 1 })
	defer close(release)

	go func() { _, _ = c.Call(context.Background(), "slow") }()
	time.Sleep(20 * time.Millisecond)

	ctx, cancel := context.WithTimeout(context.Background(), 20*time.Millisecond)
	defer cancel()
	if _, err := c.Call(ctx, "slow"); !errors.Is(err, context.DeadlineExceeded) {
		t.Fatalf("err = %v", err)
	}
}

func TestCloseRejectsFurtherCalls(t *testing.T) {
	srv := middlewaretest.NewServer(t)
	c := dial(t, srv)
	if err := c.Close(); err != nil {
		t.Fatalf("Close: %v", err)
	}
	if err := c.Close(); err != nil {
		t.Fatalf("second Close: %v", err)
	}
	if _, err := c.Call(context.Background(), "cronjob.query"); err == nil || !strings.Contains(err.Error(), "closed") {
		t.Fatalf("err = %v", err)
	}
}

type recordingLogger struct {
	mu      sync.Mutex
	entries []map[string]any
}

func (l *recordingLogger) Debug(_ context.Context, _ string, fields map[string]any) {
	l.mu.Lock()
	defer l.mu.Unlock()
	l.entries = append(l.entries, fields)
}

func TestLoggerNeverSeesParams(t *testing.T) {
	srv := middlewaretest.NewServer(t)
	srv.Handle("user.create", func(context.Context, []json.RawMessage) (any, error) { return 1, nil })
	log := &recordingLogger{}
	c := dial(t, srv, func(cfg *middleware.Config) { cfg.Logger = log })

	if _, err := c.Call(context.Background(), "user.create", map[string]any{"password": "hunter2"}); err != nil {
		t.Fatal(err)
	}
	log.mu.Lock()
	defer log.mu.Unlock()
	if len(log.entries) != 1 || log.entries[0]["method"] != "user.create" {
		t.Fatalf("entries = %+v", log.entries)
	}
	if b, _ := json.Marshal(log.entries); strings.Contains(string(b), "hunter2") {
		t.Fatal("logger received a secret")
	}
}

func TestKeepAlivePings(t *testing.T) {
	srv := middlewaretest.NewServer(t)
	srv.Handle("core.ping.echo", func(context.Context, []json.RawMessage) (any, error) { return "ok", nil })
	c := dial(t, srv, func(cfg *middleware.Config) { cfg.KeepAlive = 5 * time.Millisecond })
	time.Sleep(40 * time.Millisecond)
	if _, err := c.Call(context.Background(), "core.ping.echo"); err != nil {
		t.Fatalf("call after pings: %v", err)
	}
}

func TestTLSVerification(t *testing.T) {
	srv := middlewaretest.NewServer(t)
	srv.Handle("system.info", func(context.Context, []json.RawMessage) (any, error) { return map[string]any{}, nil })
	sum := sha256.Sum256(srv.Certificate().Raw)
	fingerprint := hex.EncodeToString(sum[:])
	caPEM := string(pem.EncodeToMemory(&pem.Block{Type: "CERTIFICATE", Bytes: srv.Certificate().Raw}))

	withTLS := func(tlsCfg middleware.TLSConfig) error {
		cfg := srv.Config()
		cfg.HTTPClient = nil
		cfg.TLS = tlsCfg
		c, err := middleware.Dial(context.Background(), cfg)
		if err == nil {
			_ = c.Close()
		}
		return err
	}

	colons := strings.ToUpper(fingerprint[:2]) + ":" + fingerprint[2:]
	for name, cfg := range map[string]middleware.TLSConfig{
		"fingerprint":        {Fingerprint: fingerprint},
		"fingerprint colons": {Fingerprint: "sha256:" + colons},
		"ca pem":             {CAPEM: caPEM},
		"insecure":           {InsecureSkipVerify: true},
	} {
		if err := withTLS(cfg); err != nil {
			t.Errorf("%s: %v", name, err)
		}
	}

	wrong := strings.Repeat("ab", sha256.Size)
	for name, cfg := range map[string]middleware.TLSConfig{
		"system roots":      {},
		"wrong fingerprint": {Fingerprint: wrong},
		"bad fingerprint":   {Fingerprint: "not-hex"},
		"bad ca":            {CAPEM: "not a pem"},
	} {
		if err := withTLS(cfg); err == nil {
			t.Errorf("%s: expected failure", name)
		}
	}
}

func TestCallIntoPropagatesCallError(t *testing.T) {
	srv := middlewaretest.NewServer(t)
	c := dial(t, srv)
	var out any
	if err := c.CallInto(context.Background(), "missing.method", &out); err == nil {
		t.Fatal("expected error")
	}
	if got := srv.Methods(); len(got) != 1 || got[0] != "missing.method" {
		t.Errorf("Methods() = %v", got)
	}
	raw, err := c.Call(context.Background(), "core.ping")
	if err != nil || string(raw) != `"pong"` {
		t.Errorf("core.ping = %s, %v", raw, err)
	}
}
