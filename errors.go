package fipe

import (
	"errors"
	"fmt"
	"net/http"
)

// Sentinel errors for common API failures. APIError unwraps to these, so
// callers can match with errors.Is regardless of which form they prefer.
var (
	ErrNotFound        = errors.New("fipe: not found")
	ErrTooManyRequests = errors.New("fipe: too many requests")
)

// APIError is returned when the API responds with a non-2xx status.
type APIError struct {
	StatusCode int
	Body       string
}

func (e *APIError) Error() string {
	if e.Body != "" {
		return fmt.Sprintf("fipe: unexpected status %d: %s", e.StatusCode, e.Body)
	}
	return fmt.Sprintf("fipe: unexpected status %d", e.StatusCode)
}

func (e *APIError) Unwrap() error {
	switch e.StatusCode {
	case http.StatusNotFound:
		return ErrNotFound
	case http.StatusTooManyRequests:
		return ErrTooManyRequests
	}
	return nil
}
