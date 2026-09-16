package middleware

import (
	"encoding/json"
	"errors"
	"fmt"
	"strings"
	"syscall"
)

// JSON-RPC error codes used by the TrueNAS middleware.
const (
	CodeInvalidJSON            = -32700
	CodeInvalidRequest         = -32600
	CodeMethodNotFound         = -32601
	CodeInvalidParams          = -32602
	CodeInternalError          = -32603
	CodeTooManyConcurrentCalls = -32000
	CodeCallError              = -32001
)

const errnameNotFound = "ENOENT"

// errnameBusy is middleware's "come back shortly": rate limits and contended resources both use
// it. It is a pacing signal, not a failure, so a caller should back off rather than give up.
const errnameBusy = "EBUSY"

// Error is an error reported by the middleware for a single call.
type Error struct {
	// Method is the API method that failed.
	Method string
	// Code is the JSON-RPC error code.
	Code int
	// Message is the JSON-RPC error message, e.g. "Method call error".
	Message string
	// Errno is the POSIX errno the middleware attached, or 0.
	Errno int
	// Errname is the symbolic errno name, e.g. "ENOENT".
	Errname string
	// Reason is the human-readable failure description.
	Reason string
	// Fields lists per-attribute validation failures, when the middleware reported any.
	Fields []FieldError
	// Trace is the formatted server-side traceback, when present.
	Trace string
}

// FieldError is one attribute-level validation failure.
type FieldError struct {
	// Attribute is the dotted path the middleware reported, e.g. "sharing_smb_create.path".
	Attribute string
	Message   string
	Errno     int
}

func (e *Error) Error() string {
	var b strings.Builder
	b.WriteString(e.Method)
	b.WriteString(": ")
	switch {
	case e.Reason != "":
		b.WriteString(e.Reason)
	case e.Message != "":
		b.WriteString(e.Message)
	default:
		fmt.Fprintf(&b, "JSON-RPC error %d", e.Code)
	}
	if e.Errname != "" {
		fmt.Fprintf(&b, " [%s]", e.Errname)
	}
	return b.String()
}

// IsNotFound reports whether err is a middleware ENOENT error.
func IsNotFound(err error) bool {
	var e *Error
	return errors.As(err, &e) && (e.Errname == errnameNotFound || e.Errno == int(syscall.ENOENT))
}

// IsValidation reports whether err carries attribute-level validation failures.
func IsValidation(err error) bool {
	var e *Error
	return errors.As(err, &e) && len(e.Fields) > 0
}

// ConnectionError reports that the websocket connection failed.
type ConnectionError struct {
	// Method is the call that was in flight, if any.
	Method string
	// Sent is true when the request may have reached the server. The call's
	// outcome is then unknown and it must not be retried blindly.
	Sent bool
	Err  error
}

func (e *ConnectionError) Error() string {
	state := "before the request was sent"
	if e.Sent {
		state = "after the request was sent; outcome unknown"
	}
	if e.Method == "" {
		return fmt.Sprintf("truenas connection failed: %v", e.Err)
	}
	return fmt.Sprintf("%s: connection lost %s: %v", e.Method, state, e.Err)
}

func (e *ConnectionError) Unwrap() error { return e.Err }

// IsRetryable reports whether a call that failed with err is known not to have
// executed and may be sent again.
func IsRetryable(err error) bool {
	var ce *ConnectionError
	if errors.As(err, &ce) {
		return !ce.Sent
	}
	var e *Error
	return errors.As(err, &e) && e.Code == CodeTooManyConcurrentCalls
}

// AuthError reports that the middleware rejected the credentials.
type AuthError struct {
	ResponseType string
}

func (e *AuthError) Error() string {
	switch e.ResponseType {
	case "AUTH_ERR":
		return "truenas authentication failed: invalid username or API key"
	case "EXPIRED":
		return "truenas authentication failed: API key is expired or revoked. TrueNAS reports both as EXPIRED, " +
			"including keys it revokes automatically after receiving them over plaintext HTTP; check the key's " +
			"revoked_reason in the TrueNAS UI"
	case "OTP_REQUIRED":
		return "truenas authentication failed: account requires a one-time password, which API keys cannot provide"
	default:
		return fmt.Sprintf("truenas authentication failed: %s", e.ResponseType)
	}
}

// wireError is the JSON-RPC error object sent by the middleware.
type wireError struct {
	Code    int             `json:"code"`
	Message string          `json:"message"`
	Data    json.RawMessage `json:"data,omitempty"`
}

type wireErrorData struct {
	Error   *int              `json:"error"`
	Errname string            `json:"errname"`
	Reason  string            `json:"reason"`
	Trace   *wireTrace        `json:"trace"`
	Extra   []json.RawMessage `json:"extra"`
}

type wireTrace struct {
	Formatted string `json:"formatted"`
}

func (w *wireError) toError(method string) *Error {
	e := &Error{Method: method, Code: w.Code, Message: w.Message}
	if len(w.Data) == 0 {
		return e
	}
	var data wireErrorData
	if json.Unmarshal(w.Data, &data) != nil {
		return e
	}
	if data.Error != nil {
		e.Errno = *data.Error
	}
	e.Errname = data.Errname
	e.Reason = data.Reason
	if data.Trace != nil {
		e.Trace = data.Trace.Formatted
	}
	e.Fields = parseFieldErrors(data.Extra)
	return e
}

// parseFieldErrors decodes validation extras of the form [attribute, message, errno].
// Extras of any other shape (CallError payloads) are ignored.
func parseFieldErrors(extra []json.RawMessage) []FieldError {
	var fields []FieldError
	for _, raw := range extra {
		var tuple []json.RawMessage
		if json.Unmarshal(raw, &tuple) != nil || len(tuple) != 3 {
			continue
		}
		var f FieldError
		if json.Unmarshal(tuple[0], &f.Attribute) != nil || json.Unmarshal(tuple[1], &f.Message) != nil {
			continue
		}
		_ = json.Unmarshal(tuple[2], &f.Errno)
		fields = append(fields, f)
	}
	return fields
}

// IsBusy reports whether err means the server is asking the caller to slow down: a rate limit, or
// a resource briefly in use. Such a call is worth retrying after a wait.
func IsBusy(err error) bool {
	var e *Error
	if !errors.As(err, &e) {
		return false
	}
	return e.Errname == errnameBusy || e.Errno == int(syscall.EBUSY) || e.Code == CodeTooManyConcurrentCalls
}
