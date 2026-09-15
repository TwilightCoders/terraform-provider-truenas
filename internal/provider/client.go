package provider

import (
	"errors"

	"github.com/TwilightCoders/terraform-provider-truenas/internal/engine"
	"github.com/TwilightCoders/terraform-provider-truenas/internal/middleware"
)

// engineClient adapts the middleware client to what the engine asks of a transport.
//
// The engine classifies errors by asking rather than by inspecting types, so it carries no
// dependency on this provider's transport. The adapter lives here, on the provider side, so that
// neither the engine nor the middleware package depends on the other.
type engineClient struct{ *middleware.Client }

var _ engine.FileClient = engineClient{}

// IsNotFound reports whether err means the object does not exist.
func (c engineClient) IsNotFound(err error) bool { return middleware.IsNotFound(err) }

// IsRetryable reports whether err is a dropped connection. Only calls the engine knows to be
// idempotent are retried on it.
func (c engineClient) IsRetryable(err error) bool {
	var ce *middleware.ConnectionError
	return errors.As(err, &ce)
}

// ValidationFields returns the per-attribute validation failures err carries, if any.
func (c engineClient) ValidationFields(err error) []engine.FieldError {
	var e *middleware.Error
	if !errors.As(err, &e) || len(e.Fields) == 0 {
		return nil
	}
	out := make([]engine.FieldError, 0, len(e.Fields))
	for _, f := range e.Fields {
		out = append(out, engine.FieldError{Attribute: f.Attribute, Message: f.Message})
	}
	return out
}
