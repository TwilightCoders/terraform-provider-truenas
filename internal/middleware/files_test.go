package middleware_test

import (
	"bytes"
	"context"
	"encoding/json"
	"errors"
	"strings"
	"testing"

	"github.com/TwilightCoders/terraform-provider-truenas/internal/middleware"
	"github.com/TwilightCoders/terraform-provider-truenas/internal/middleware/middlewaretest"
)

func TestGetFile(t *testing.T) {
	s := middlewaretest.NewServer(t)
	s.Files.Put("/mnt/pool/app/compose.yml", []byte("services: {}\n"))
	c := dial(t, s)

	var buf bytes.Buffer
	if err := c.GetFile(context.Background(), "/mnt/pool/app/compose.yml", &buf); err != nil {
		t.Fatalf("GetFile: %v", err)
	}
	if got := buf.String(); got != "services: {}\n" {
		t.Errorf("GetFile = %q", got)
	}
}

func TestGetFileExpiredJob(t *testing.T) {
	s := middlewaretest.NewServer(t)
	s.Files.Put("/mnt/pool/f", []byte("x"))
	s.Files.ExpireDownloads = true
	c := dial(t, s)

	err := c.GetFile(context.Background(), "/mnt/pool/f", &bytes.Buffer{})
	if err == nil || !strings.Contains(err.Error(), "expired") {
		t.Fatalf("expected an expiry error, got %v", err)
	}
}

func TestPutFile(t *testing.T) {
	s := middlewaretest.NewServer(t)
	c := dial(t, s)

	mode := int64(0o644)
	if err := c.PutFile(context.Background(), "/mnt/pool/new.yml", &mode, strings.NewReader("hello")); err != nil {
		t.Fatalf("PutFile: %v", err)
	}
	got, ok := s.Files.Get("/mnt/pool/new.yml")
	if !ok || string(got) != "hello" {
		t.Errorf("uploaded content = %q, present=%v", got, ok)
	}
}

// TestPutFileNeverSendsAPIKey is the point of the token dance. /_upload authenticates the HTTP
// request on its own, and TrueNAS permanently revokes an API key that reaches it over a transport
// it considers insecure — which is how the first key on the house NAS died. The provider must
// therefore mint a short-lived token over the authenticated websocket and send that instead. The
// fake refuses an Authorization: Bearer header outright, so this fails if that ever regresses.
func TestPutFileNeverSendsAPIKey(t *testing.T) {
	s := middlewaretest.NewServer(t)
	c := dial(t, s)

	if err := c.PutFile(context.Background(), "/mnt/pool/f", nil, strings.NewReader("x")); err != nil {
		t.Fatalf("PutFile: %v", err)
	}
	var minted bool
	for _, call := range s.Calls() {
		if call.Method == "auth.generate_token" {
			minted = true
		}
	}
	if !minted {
		t.Error("upload did not mint a token; the API key would have been sent over HTTP")
	}
}

// TestFileTransferRefusesReverseProxy covers the guard that does not come for free: the websocket
// handshake check does not apply to the file endpoints, and those are the ones carrying credentials.
func TestFileTransferRefusesReverseProxy(t *testing.T) {
	s := middlewaretest.NewServer(t)
	s.Files.Put("/mnt/pool/f", []byte("x"))
	c := dial(t, s)

	// The proxy appears only after the session is established, as it would when a proxy fronts
	// the file endpoints but not the API path.
	// What a real reverse proxy in front of TrueNAS answers with; TrueNAS's own web server
	// answers "nginx".
	s.SetServerHeader("openresty")

	var proxyErr *middleware.ProxyError
	if err := c.GetFile(context.Background(), "/mnt/pool/f", &bytes.Buffer{}); !errors.As(err, &proxyErr) {
		t.Errorf("GetFile through a proxy = %v, want ProxyError", err)
	}
	if err := c.PutFile(context.Background(), "/mnt/pool/f", nil, strings.NewReader("x")); !errors.As(err, &proxyErr) {
		t.Errorf("PutFile through a proxy = %v, want ProxyError", err)
	}
}

func TestFileErrorsNeverLeakTheDownloadToken(t *testing.T) {
	s := middlewaretest.NewServer(t)
	s.Files.ExpireDownloads = true
	s.Files.Put("/mnt/pool/f", []byte("x"))
	c := dial(t, s)

	err := c.GetFile(context.Background(), "/mnt/pool/f", &bytes.Buffer{})
	if err == nil {
		t.Fatal("expected an error")
	}
	if strings.Contains(err.Error(), "auth_token") || strings.Contains(err.Error(), "dl-") {
		t.Errorf("error leaked the download token: %v", err)
	}
}

// TestPutFileReportsFailedJob covers the half of the job that the HTTP status cannot: /_upload
// accepts the upload before the write runs, so a write that then fails is only visible through
// the job.
func TestPutFileReportsFailedJob(t *testing.T) {
	s := middlewaretest.NewServer(t)
	s.Handle("core.job_wait", func(context.Context, []json.RawMessage) (any, error) {
		return nil, &middleware.Error{Errno: 28, Errname: "ENOSPC", Reason: "No space left on device"}
	})
	c := dial(t, s)

	err := c.PutFile(context.Background(), "/mnt/pool/f", nil, strings.NewReader("x"))
	if err == nil || !strings.Contains(err.Error(), "No space left") {
		t.Fatalf("PutFile = %v, want the job's failure", err)
	}
	if _, ok := s.Files.Get("/mnt/pool/f"); ok {
		t.Error("a failed job left a file behind")
	}
}
