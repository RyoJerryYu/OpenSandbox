package errors

import (
	"fmt"
	"time"
)

// SandboxError is the normalized SDK error surface for HTTP-backed failures.
type SandboxError struct {
	Code       string
	Message    string
	RequestID  string
	StatusCode int
	Cause      error
}

// Error implements the error interface.
func (e *SandboxError) Error() string {
	if e == nil {
		return "<nil>"
	}
	if e.Code != "" {
		return fmt.Sprintf("%s: %s", e.Code, e.Message)
	}
	return e.Message
}

// Unwrap returns the underlying cause when present.
func (e *SandboxError) Unwrap() error {
	if e == nil {
		return nil
	}
	return e.Cause
}

// ReadyTimeoutError indicates that readiness polling did not succeed before the deadline.
type ReadyTimeoutError struct {
	Timeout time.Duration
}

// Error implements the error interface.
func (e *ReadyTimeoutError) Error() string {
	if e == nil {
		return "<nil>"
	}
	return fmt.Sprintf("sandbox was not ready within %s", e.Timeout)
}

// InvalidArgumentError reports invalid user input before an HTTP request is made.
type InvalidArgumentError struct {
	Message string
}

// Error implements the error interface.
func (e *InvalidArgumentError) Error() string {
	if e == nil {
		return "<nil>"
	}
	return e.Message
}
