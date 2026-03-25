package unit

import (
	"errors"
	"testing"
	"time"

	sandboxerrors "github.com/alibaba/opensandbox/sdks/sandbox/go/sandbox/errors"
)

func TestSandboxErrorCarriesRequestID(t *testing.T) {
	err := &sandboxerrors.SandboxError{
		Code:      "bad_request",
		Message:   "failed",
		RequestID: "req-123",
	}

	if err.RequestID != "req-123" {
		t.Fatalf("unexpected request id: %s", err.RequestID)
	}
}

func TestReadyTimeoutErrorMatchesErrorsAs(t *testing.T) {
	target := &sandboxerrors.ReadyTimeoutError{}
	err := &sandboxerrors.ReadyTimeoutError{Timeout: 30 * time.Second}

	if !errors.As(err, &target) {
		t.Fatal("expected errors.As to match ReadyTimeoutError")
	}
}
