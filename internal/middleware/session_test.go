package middleware

import (
	"context"
	"encoding/json"
	"errors"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"
	"time"

	"github.com/coder/websocket"
)

// scriptedServer answers each request with reply(method, id), letting tests send protocol
// messages the fake middleware never produces.
func scriptedServer(t *testing.T, reply func(ctx context.Context, conn *websocket.Conn, method string, id json.RawMessage)) Config {
	t.Helper()
	srv := httptest.NewTLSServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.URL.Path == "/api/versions" {
			_, _ = w.Write([]byte(`["v25.10.5"]`))
			return
		}
		conn, err := websocket.Accept(w, r, nil)
		if err != nil {
			return
		}
		defer func() { _ = conn.CloseNow() }()
		for {
			_, data, err := conn.Read(r.Context())
			if err != nil {
				return
			}
			var req struct {
				ID     json.RawMessage `json:"id"`
				Method string          `json:"method"`
			}
			_ = json.Unmarshal(data, &req)
			reply(r.Context(), conn, req.Method, req.ID)
		}
	}))
	t.Cleanup(srv.Close)
	return Config{Host: srv.Listener.Addr().String(), APIVersion: "v25.10.5", Username: "u", APIKey: "k", HTTPClient: srv.Client()}
}

func write(ctx context.Context, conn *websocket.Conn, s string) {
	_ = conn.Write(ctx, websocket.MessageText, []byte(s))
}

func okSetup(ctx context.Context, conn *websocket.Conn, method string, id json.RawMessage) bool {
	switch method {
	case "core.set_options":
		write(ctx, conn, `{"jsonrpc":"2.0","id":`+string(id)+`,"result":{}}`)
		return true
	case "auth.login_ex":
		write(ctx, conn, `{"jsonrpc":"2.0","id":`+string(id)+`,"result":{"response_type":"SUCCESS"}}`)
		return true
	}
	return false
}

func TestSessionIgnoresNoiseAndDeliversResponse(t *testing.T) {
	cfg := scriptedServer(t, func(ctx context.Context, conn *websocket.Conn, method string, id json.RawMessage) {
		if okSetup(ctx, conn, method, id) {
			return
		}
		write(ctx, conn, `not json`)
		write(ctx, conn, `{"jsonrpc":"2.0","method":"collection_update","params":{"msg":"added"}}`)
		write(ctx, conn, `{"jsonrpc":"2.0","id":999999,"result":"stray"}`)
		write(ctx, conn, `{"jsonrpc":"2.0","id":`+string(id)+`,"result":"ok"}`)
	})
	c, err := Dial(context.Background(), cfg)
	if err != nil {
		t.Fatal(err)
	}
	defer func() { _ = c.Close() }()
	raw, err := c.Call(context.Background(), "x")
	if err != nil || string(raw) != `"ok"` {
		t.Fatalf("got %s, %v", raw, err)
	}
}

func TestSessionGlobalErrorKillsSession(t *testing.T) {
	cfg := scriptedServer(t, func(ctx context.Context, conn *websocket.Conn, method string, id json.RawMessage) {
		if okSetup(ctx, conn, method, id) {
			return
		}
		write(ctx, conn, `{"jsonrpc":"2.0","id":null,"error":{"code":-32700,"message":"Invalid JSON"}}`)
	})
	c, err := Dial(context.Background(), cfg)
	if err != nil {
		t.Fatal(err)
	}
	defer func() { _ = c.Close() }()
	_, err = c.Call(context.Background(), "x")
	var ce *ConnectionError
	if !errors.As(err, &ce) || !ce.Sent {
		t.Fatalf("err = %v", err)
	}
}

func TestConnectSetupFailures(t *testing.T) {
	tests := map[string]struct {
		reply   func(ctx context.Context, conn *websocket.Conn, method string, id json.RawMessage)
		errPart string
	}{
		"set_options error": {
			reply: func(ctx context.Context, conn *websocket.Conn, _ string, id json.RawMessage) {
				write(ctx, conn, `{"jsonrpc":"2.0","id":`+string(id)+`,"error":{"code":-32601,"message":"Method does not exist"}}`)
			},
			errPart: "core.set_options",
		},
		"login malformed": {
			reply: func(ctx context.Context, conn *websocket.Conn, method string, id json.RawMessage) {
				if method == "core.set_options" {
					okSetup(ctx, conn, method, id)
					return
				}
				write(ctx, conn, `{"jsonrpc":"2.0","id":`+string(id)+`,"result":"yes"}`)
			},
			errPart: "authenticating",
		},
	}
	for name, tt := range tests {
		t.Run(name, func(t *testing.T) {
			_, err := Dial(context.Background(), scriptedServer(t, tt.reply))
			if err == nil || !strings.Contains(err.Error(), tt.errPart) {
				t.Fatalf("err = %v, want containing %q", err, tt.errPart)
			}
		})
	}
}

func TestSessionCallMarshalError(t *testing.T) {
	cfg := scriptedServer(t, func(ctx context.Context, conn *websocket.Conn, method string, id json.RawMessage) {
		okSetup(ctx, conn, method, id)
	})
	c, err := Dial(context.Background(), cfg)
	if err != nil {
		t.Fatal(err)
	}
	defer func() { _ = c.Close() }()
	if _, err := c.Call(context.Background(), "x", make(chan int)); err == nil {
		t.Fatal("expected marshal error")
	}
}

func TestSessionFailsFastAfterTeardown(t *testing.T) {
	cfg := scriptedServer(t, func(ctx context.Context, conn *websocket.Conn, method string, id json.RawMessage) {
		okSetup(ctx, conn, method, id)
	})
	c, err := Dial(context.Background(), cfg)
	if err != nil {
		t.Fatal(err)
	}
	defer func() { _ = c.Close() }()

	s := c.sess
	s.fail(errors.New("synthetic"))
	if _, err := s.call(context.Background(), 1, "x", nil); err == nil || !strings.Contains(err.Error(), "synthetic") {
		t.Fatalf("call on dead session: %v", err)
	}
	if err := s.conn.Write(context.Background(), websocket.MessageText, []byte("{}")); err == nil {
		t.Fatal("expected write on closed conn to fail")
	}
}

func TestWriteFailureMarksCallSent(t *testing.T) {
	cfg := scriptedServer(t, func(ctx context.Context, conn *websocket.Conn, method string, id json.RawMessage) {
		okSetup(ctx, conn, method, id)
	})
	c, err := Dial(context.Background(), cfg)
	if err != nil {
		t.Fatal(err)
	}
	defer func() { _ = c.Close() }()

	s := c.sess
	_ = s.conn.CloseNow() // the reader has not noticed yet when the write fails
	_, err = s.call(context.Background(), 42, "x", nil)
	var ce *ConnectionError
	if !errors.As(err, &ce) {
		t.Fatalf("err = %v", err)
	}
}

func TestPingFailureKillsSession(t *testing.T) {
	blocked := make(chan struct{})
	t.Cleanup(func() { close(blocked) })
	srv := httptest.NewTLSServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.URL.Path == "/api/versions" {
			_, _ = w.Write([]byte(`["v25.10.5"]`))
			return
		}
		conn, err := websocket.Accept(w, r, nil)
		if err != nil {
			return
		}
		ctx := r.Context()
		for range 2 {
			_, data, err := conn.Read(ctx)
			if err != nil {
				return
			}
			var req struct {
				ID     json.RawMessage `json:"id"`
				Method string          `json:"method"`
			}
			_ = json.Unmarshal(data, &req)
			okSetup(ctx, conn, req.Method, req.ID)
		}
		<-blocked // stop reading: pings are never answered
	}))
	t.Cleanup(srv.Close)

	cfg := Config{Host: srv.Listener.Addr().String(), APIVersion: "v25.10.5", Username: "u", APIKey: "k", HTTPClient: srv.Client(), KeepAlive: 10 * time.Millisecond}
	c, err := Dial(context.Background(), cfg)
	if err != nil {
		t.Fatal(err)
	}
	defer func() { _ = c.Close() }()

	deadline := time.After(2 * time.Second)
	for c.sess.alive() {
		select {
		case <-deadline:
			t.Fatal("session survived unanswered pings")
		case <-time.After(10 * time.Millisecond):
		}
	}
}
