package middleware

import (
	"crypto/tls"
	"encoding/json"
	"errors"
	"strings"
	"testing"
)

func TestSelectVersion(t *testing.T) {
	tests := []struct {
		name      string
		available []string
		want      string
		expect    string
		errPart   string
	}{
		{name: "exact", available: []string{"v25.10.5"}, want: "v25.10.5", expect: "v25.10.5"},
		{name: "older patch", available: []string{"v25.04.2", "v25.10.1", "v25.10.2"}, want: "v25.10.5", expect: "v25.10.2"},
		{name: "ignores newer patch", available: []string{"v25.10.4", "v25.10.9"}, want: "v25.10.5", expect: "v25.10.4"},
		{name: "ignores other series", available: []string{"v26.0.0", "v25.04.2"}, want: "v25.10.5", errPart: "needs v25.10.x"},
		{name: "ignores junk", available: []string{"current", "v25.10", "v25.10.x", "25.10.1", "v25.10.-1"}, want: "v25.10.5", errPart: "needs"},
		{name: "invalid want", available: []string{"v25.10.5"}, want: "25.10", errPart: "invalid API version"},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := SelectVersion(tt.available, tt.want)
			if tt.errPart != "" {
				if err == nil || !strings.Contains(err.Error(), tt.errPart) {
					t.Fatalf("err = %v, want containing %q", err, tt.errPart)
				}
				return
			}
			if err != nil || got != tt.expect {
				t.Fatalf("got %q, %v; want %q", got, err, tt.expect)
			}
		})
	}
}

func TestParseFieldErrorsSkipsNonValidationExtras(t *testing.T) {
	raw := []json.RawMessage{
		json.RawMessage(`["user_create.username", "already exists", 17]`),
		json.RawMessage(`{"app": "plex"}`),
		json.RawMessage(`["too", "short"]`),
		json.RawMessage(`[1, "attribute not a string", 22]`),
		json.RawMessage(`["x", 2, 22]`),
	}
	got := parseFieldErrors(raw)
	if len(got) != 1 || got[0] != (FieldError{Attribute: "user_create.username", Message: "already exists", Errno: 17}) {
		t.Fatalf("got %+v", got)
	}
}

func TestWireErrorToleratesMissingOrOddData(t *testing.T) {
	e := (&wireError{Code: CodeInvalidRequest, Message: "bad"}).toError("m")
	if e.Error() != "m: bad" {
		t.Errorf("Error() = %q", e.Error())
	}
	e = (&wireError{Code: CodeInternalError, Data: json.RawMessage(`"not an object"`)}).toError("m")
	if e.Error() != "m: JSON-RPC error -32603" {
		t.Errorf("Error() = %q", e.Error())
	}
	e = (&wireError{Code: CodeCallError, Data: json.RawMessage(`{"error":2,"reason":"gone","trace":{"formatted":"Traceback"}}`)}).toError("m")
	if e.Errno != 2 || e.Trace != "Traceback" || !IsNotFound(e) {
		t.Errorf("got %#v", e)
	}
}

func TestConnectionErrorWithoutMethod(t *testing.T) {
	err := &ConnectionError{Err: errors.New("dial tcp: refused")}
	if err.Error() != "truenas connection failed: dial tcp: refused" {
		t.Errorf("Error() = %q", err.Error())
	}
	if !errors.Is(err, err.Err) {
		t.Error("Unwrap broken")
	}
}

func TestTLSBuildDefaults(t *testing.T) {
	cfg, err := TLSConfig{ServerName: "nas.lan"}.build()
	if err != nil {
		t.Fatal(err)
	}
	if cfg.MinVersion != tls.VersionTLS12 || cfg.ServerName != "nas.lan" || cfg.InsecureSkipVerify {
		t.Errorf("unexpected config %+v", cfg)
	}
}

func TestPinnedFingerprintRejectsMissingCertificate(t *testing.T) {
	cfg, err := TLSConfig{Fingerprint: strings.Repeat("00", 32)}.build()
	if err != nil {
		t.Fatal(err)
	}
	if err := cfg.VerifyConnection(tls.ConnectionState{}); err == nil {
		t.Fatal("expected error for missing peer certificate")
	}
}
