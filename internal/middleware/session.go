package middleware

import (
	"bytes"
	"context"
	"encoding/json"
	"errors"
	"sync"
	"time"

	"github.com/coder/websocket"
)

// session is one authenticated websocket connection. It dies on the first transport error.
type session struct {
	conn   *websocket.Conn
	ctx    context.Context
	cancel context.CancelFunc
	done   chan struct{}

	mu      sync.Mutex
	pending map[uint64]chan response
	err     error
}

type request struct {
	JSONRPC string `json:"jsonrpc"`
	ID      uint64 `json:"id"`
	Method  string `json:"method"`
	Params  []any  `json:"params"`
}

type response struct {
	ID     *uint64         `json:"id"`
	Method string          `json:"method"`
	Result json.RawMessage `json:"result"`
	Error  *wireError      `json:"error"`
}

func newSession(conn *websocket.Conn, keepAlive time.Duration) *session {
	ctx, cancel := context.WithCancel(context.Background())
	s := &session{
		conn:    conn,
		ctx:     ctx,
		cancel:  cancel,
		done:    make(chan struct{}),
		pending: make(map[uint64]chan response),
	}
	go s.readLoop()
	if keepAlive > 0 {
		go s.pingLoop(keepAlive)
	}
	return s
}

func (s *session) alive() bool {
	select {
	case <-s.done:
		return false
	default:
		return true
	}
}

func (s *session) call(ctx context.Context, id uint64, method string, params []any) (json.RawMessage, error) {
	if params == nil {
		params = []any{}
	}
	payload, err := json.Marshal(request{JSONRPC: "2.0", ID: id, Method: method, Params: params})
	if err != nil {
		return nil, err
	}

	ch := make(chan response, 1)
	s.mu.Lock()
	if s.err != nil {
		err := s.err
		s.mu.Unlock()
		return nil, &ConnectionError{Method: method, Err: err}
	}
	s.pending[id] = ch
	s.mu.Unlock()
	defer s.unregister(id)

	// The write gets its own deadline: cancelling a caller's context mid-write would close the
	// shared connection under every other in-flight call.
	wctx, cancel := context.WithTimeout(s.ctx, defaultWriteTimeout)
	err = s.conn.Write(wctx, websocket.MessageText, payload)
	cancel()
	if err != nil {
		s.fail(err)
		return nil, &ConnectionError{Method: method, Sent: true, Err: err}
	}

	select {
	case resp := <-ch:
		return resp.unwrap(method)
	case <-s.done:
		// A response that raced the teardown still wins.
		select {
		case resp := <-ch:
			return resp.unwrap(method)
		default:
			return nil, &ConnectionError{Method: method, Sent: true, Err: s.failure()}
		}
	case <-ctx.Done():
		return nil, ctx.Err()
	}
}

func (r response) unwrap(method string) (json.RawMessage, error) {
	if r.Error != nil {
		return nil, r.Error.toError(method)
	}
	return r.Result, nil
}

func (s *session) unregister(id uint64) {
	s.mu.Lock()
	delete(s.pending, id)
	s.mu.Unlock()
}

func (s *session) readLoop() {
	for {
		_, data, err := s.conn.Read(s.ctx)
		if err != nil {
			s.fail(err)
			return
		}
		var resp response
		dec := json.NewDecoder(bytes.NewReader(data))
		if err := dec.Decode(&resp); err != nil {
			continue
		}
		if resp.ID == nil {
			if resp.Method == "" && resp.Error != nil {
				// A global error with a null id terminates the session.
				s.fail(resp.Error.toError(""))
				return
			}
			continue // notification; this client does not subscribe to events
		}
		s.mu.Lock()
		ch, ok := s.pending[*resp.ID]
		s.mu.Unlock()
		if ok {
			ch <- resp
		}
	}
}

func (s *session) pingLoop(interval time.Duration) {
	t := time.NewTicker(interval)
	defer t.Stop()
	for {
		select {
		case <-s.done:
			return
		case <-t.C:
			ctx, cancel := context.WithTimeout(s.ctx, interval)
			err := s.conn.Ping(ctx)
			cancel()
			if err != nil {
				s.fail(err)
				return
			}
		}
	}
}

// fail records the first transport error and tears the session down.
func (s *session) fail(err error) {
	s.mu.Lock()
	if s.err != nil {
		s.mu.Unlock()
		return
	}
	s.err = err
	s.mu.Unlock()
	s.cancel()
	_ = s.conn.CloseNow()
	close(s.done)
}

func (s *session) failure() error {
	s.mu.Lock()
	defer s.mu.Unlock()
	return s.err
}

func (s *session) close() error {
	err := s.conn.Close(websocket.StatusNormalClosure, "")
	s.fail(errors.New("client closed"))
	return err
}
