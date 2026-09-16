// Package middlewaretest provides an in-process TrueNAS middleware for tests.
//
// Server speaks the same JSON-RPC 2.0 websocket protocol as TrueNAS: version discovery at
// /api/versions, a websocket per API version, auth.login_ex and core.set_options. Tests register
// handlers for the methods they exercise.
package middlewaretest

import (
	"context"
	"encoding/json"
	"errors"
	"net/http"
	"net/http/httptest"
	"slices"
	"strings"
	"sync"
	"testing"

	"github.com/coder/websocket"

	"github.com/TwilightCoders/terraform-provider-truenas/internal/middleware"
)

// Default credentials and API version served by NewServer.
const (
	Username   = "root"
	APIKey     = "1-test-api-key"
	APIVersion = "v25.10.5"
)

// Handler serves one API method. Returning a *middleware.Error sends that error to the client;
// any other error is sent as a generic call error.
type Handler func(ctx context.Context, params []json.RawMessage) (any, error)

// Call records one method invocation received by the server.
type Call struct {
	Method string
	Params []json.RawMessage
}

// Server is a fake TrueNAS middleware.
type Server struct {
	*httptest.Server

	mu            sync.Mutex
	versions      []string
	handlers      map[string]Handler
	calls         []Call
	conns         map[*websocket.Conn]struct{}
	loginResponse string
	loginHook     func() error
	serverHeader  string

	// Files backs the /_download and /_upload endpoints.
	Files *Files
}

// NewServer starts a TLS server offering APIVersion. It is closed when the test ends.
func NewServer(t testing.TB) *Server {
	t.Helper()
	return newServer(t, httptest.NewTLSServer)
}

// NewPlaintextServer starts a plaintext server on 127.0.0.1, like the local end of an SSH tunnel.
func NewPlaintextServer(t testing.TB) *Server {
	t.Helper()
	return newServer(t, httptest.NewServer)
}

func newServer(t testing.TB, start func(http.Handler) *httptest.Server) *Server {
	t.Helper()
	s := &Server{
		versions:      []string{APIVersion},
		handlers:      make(map[string]Handler),
		conns:         make(map[*websocket.Conn]struct{}),
		loginResponse: "SUCCESS",
		serverHeader:  "nginx",
		Files:         newFiles(),
	}
	mux := http.NewServeMux()
	mux.HandleFunc("GET /api/versions", s.serveVersions)
	mux.HandleFunc("/api/{version}", s.serveWebsocket)
	mux.HandleFunc("/_download/", s.serveDownload)
	mux.HandleFunc("POST /_upload", s.serveUpload)
	s.Server = start(mux)
	t.Cleanup(s.Close)
	return s
}

// Config returns a client config that trusts and authenticates to this server.
func (s *Server) Config() middleware.Config {
	return middleware.Config{
		Host:       s.Listener.Addr().String(),
		APIVersion: APIVersion,
		Username:   Username,
		APIKey:     APIKey,
		HTTPClient: s.Client(),
	}
}

// SetVersions replaces the API versions the server offers.
func (s *Server) SetVersions(versions ...string) {
	s.mu.Lock()
	defer s.mu.Unlock()
	s.versions = versions
}

// SetLoginResponse makes auth.login_ex answer with the given response_type.
func (s *Server) SetLoginResponse(responseType string) {
	s.mu.Lock()
	defer s.mu.Unlock()
	s.loginResponse = responseType
}

// SetLoginHook runs before every auth.login_ex and, when it returns an error, fails that attempt.
// Tests use it to make the server rate limit a login the way middleware does.
func (s *Server) SetLoginHook(fn func() error) {
	s.mu.Lock()
	defer s.mu.Unlock()
	s.loginHook = fn
}

// Handle registers h for method, replacing any previous handler.
func (s *Server) Handle(method string, h Handler) {
	s.mu.Lock()
	defer s.mu.Unlock()
	s.handlers[method] = h
}

// Calls returns every call received so far, in arrival order, excluding session setup.
func (s *Server) Calls() []Call {
	s.mu.Lock()
	defer s.mu.Unlock()
	return slices.Clone(s.calls)
}

// Methods returns the method names of Calls.
func (s *Server) Methods() []string {
	calls := s.Calls()
	names := make([]string, len(calls))
	for i, c := range calls {
		names[i] = c.Method
	}
	return names
}

// DropConnections abruptly closes every open websocket.
func (s *Server) DropConnections() {
	s.mu.Lock()
	conns := make([]*websocket.Conn, 0, len(s.conns))
	for c := range s.conns {
		conns = append(conns, c)
	}
	s.mu.Unlock()
	for _, c := range conns {
		_ = c.CloseNow()
	}
}

// SetServerHeader changes the HTTP Server header, to simulate a reverse proxy.
func (s *Server) SetServerHeader(value string) {
	s.mu.Lock()
	defer s.mu.Unlock()
	s.serverHeader = value
}

// setServerHeader writes the Server header the reverse-proxy guard inspects.
func (s *Server) setServerHeader(w http.ResponseWriter) {
	s.mu.Lock()
	header := s.serverHeader
	s.mu.Unlock()
	w.Header().Set("Server", header)
}

func (s *Server) serveVersions(w http.ResponseWriter, _ *http.Request) {
	s.mu.Lock()
	versions := slices.Clone(s.versions)
	header := s.serverHeader
	s.mu.Unlock()
	w.Header().Set("Server", header)
	w.Header().Set("Content-Type", "application/json")
	_ = json.NewEncoder(w).Encode(versions)
}

func (s *Server) serveWebsocket(w http.ResponseWriter, r *http.Request) {
	s.mu.Lock()
	offered := slices.Contains(s.versions, r.PathValue("version"))
	s.mu.Unlock()
	if !offered {
		http.NotFound(w, r)
		return
	}
	conn, err := websocket.Accept(w, r, nil)
	if err != nil {
		return
	}
	conn.SetReadLimit(64 << 20)
	s.track(conn, true)
	defer s.track(conn, false)

	sess := &serverSession{server: s, conn: conn}
	ctx := r.Context()
	for {
		_, data, err := conn.Read(ctx)
		if err != nil {
			return
		}
		var req struct {
			ID     json.RawMessage   `json:"id"`
			Method string            `json:"method"`
			Params []json.RawMessage `json:"params"`
		}
		if err := json.Unmarshal(data, &req); err != nil {
			sess.send(ctx, map[string]any{"jsonrpc": "2.0", "id": nil, "error": map[string]any{"code": middleware.CodeInvalidJSON, "message": err.Error()}})
			continue
		}
		go sess.dispatch(ctx, req.ID, req.Method, req.Params)
	}
}

func (s *Server) track(c *websocket.Conn, open bool) {
	s.mu.Lock()
	defer s.mu.Unlock()
	if open {
		s.conns[c] = struct{}{}
	} else {
		delete(s.conns, c)
	}
}

type serverSession struct {
	server *Server
	conn   *websocket.Conn

	mu            sync.Mutex
	authenticated bool
}

func (ss *serverSession) dispatch(ctx context.Context, id json.RawMessage, method string, params []json.RawMessage) {
	result, err := ss.handle(ctx, method, params)
	msg := map[string]any{"jsonrpc": "2.0", "id": id}
	if err != nil {
		msg["error"] = wireError(err)
	} else {
		msg["result"] = result
	}
	ss.send(ctx, msg)
}

func (ss *serverSession) handle(ctx context.Context, method string, params []json.RawMessage) (any, error) {
	switch method {
	case "core.set_options":
		return map[string]any{"legacy_jobs": false, "private_methods": false, "py_exceptions": false}, nil
	case "core.ping":
		return "pong", nil
	case "auth.login_ex":
		return ss.login(params)
	}

	ss.mu.Lock()
	authed := ss.authenticated
	ss.mu.Unlock()
	if !authed {
		return nil, &middleware.Error{Code: middleware.CodeCallError, Errno: 207, Errname: "ENOTAUTHENTICATED", Reason: "Not authenticated"}
	}

	s := ss.server
	s.mu.Lock()
	h, ok := s.handlers[method]
	s.calls = append(s.calls, Call{Method: method, Params: params})
	s.mu.Unlock()
	if !ok {
		if result, handled, err := s.handleFileMethod(method, params); handled {
			return result, err
		}
	}
	if !ok {
		return nil, &middleware.Error{Code: middleware.CodeMethodNotFound, Message: "Method does not exist"}
	}
	return h(ctx, params)
}

func (ss *serverSession) login(params []json.RawMessage) (any, error) {
	var creds struct {
		Mechanism string `json:"mechanism"`
		Username  string `json:"username"`
		APIKey    string `json:"api_key"`
	}
	if len(params) != 1 || json.Unmarshal(params[0], &creds) != nil || creds.Mechanism != "API_KEY_PLAIN" {
		return nil, &middleware.Error{Code: middleware.CodeInvalidParams, Errname: "EINVAL", Reason: "invalid login_data"}
	}
	s := ss.server
	s.mu.Lock()
	hook := s.loginHook
	s.mu.Unlock()
	if hook != nil {
		if err := hook(); err != nil {
			return nil, err
		}
	}
	s.mu.Lock()
	responseType := s.loginResponse
	s.mu.Unlock()
	if responseType == "SUCCESS" && (creds.Username != Username || creds.APIKey != APIKey) {
		responseType = "AUTH_ERR"
	}
	if responseType == "SUCCESS" {
		ss.mu.Lock()
		ss.authenticated = true
		ss.mu.Unlock()
	}
	return map[string]any{"response_type": responseType}, nil
}

func (ss *serverSession) send(ctx context.Context, msg any) {
	data, err := json.Marshal(msg)
	if err != nil {
		panic(err)
	}
	_ = ss.conn.Write(ctx, websocket.MessageText, data)
}

// wireError renders err the way middlewared does.
func wireError(err error) map[string]any {
	var e *middleware.Error
	if !errors.As(err, &e) {
		e = &middleware.Error{Code: middleware.CodeCallError, Errno: 22, Errname: "EINVAL", Reason: err.Error()}
	}
	code, message := e.Code, e.Message
	if code == 0 {
		code = middleware.CodeCallError
	}
	if message == "" {
		message = "Method call error"
		if len(e.Fields) > 0 {
			code, message = middleware.CodeInvalidParams, "Invalid params"
		}
	}
	out := map[string]any{"code": code, "message": message}
	if e.Errname == "" && e.Reason == "" && len(e.Fields) == 0 {
		return out
	}
	extra := make([]any, 0, len(e.Fields))
	for _, f := range e.Fields {
		extra = append(extra, []any{f.Attribute, f.Message, f.Errno})
	}
	reason := e.Reason
	if reason == "" && len(e.Fields) > 0 {
		lines := make([]string, len(e.Fields))
		for i, f := range e.Fields {
			lines[i] = "[EINVAL] " + f.Attribute + ": " + f.Message
		}
		reason = strings.Join(lines, "\n")
	}
	out["data"] = map[string]any{
		"error":   e.Errno,
		"errname": e.Errname,
		"reason":  reason,
		"trace":   nil,
		"extra":   extra,
	}
	return out
}
