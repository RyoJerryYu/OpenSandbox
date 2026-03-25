package errors

import (
	"fmt"
	"time"
)

type SandboxError struct {
	Code       string
	Message    string
	RequestID  string
	StatusCode int
	Cause      error
}

func (e *SandboxError) Error() string {
	if e == nil {
		return "<nil>"
	}
	if e.Code != "" {
		return fmt.Sprintf("%s: %s", e.Code, e.Message)
	}
	return e.Message
}

func (e *SandboxError) Unwrap() error {
	if e == nil {
		return nil
	}
	return e.Cause
}

type ReadyTimeoutError struct {
	Timeout time.Duration
}

func (e *ReadyTimeoutError) Error() string {
	if e == nil {
		return "<nil>"
	}
	return fmt.Sprintf("sandbox was not ready within %s", e.Timeout)
}

type InvalidArgumentError struct {
	Message string
}

func (e *InvalidArgumentError) Error() string {
	if e == nil {
		return "<nil>"
	}
	return e.Message
}
