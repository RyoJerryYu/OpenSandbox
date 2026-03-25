package unit

import (
	"errors"
	"net/http"
	"testing"
	"time"

	"github.com/alibaba/opensandbox/sdks/sandbox/go/sandbox/adapters"
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

func TestNormalizeHTTPErrorPromotesRequestIDAndStatusCode(t *testing.T) {
	resp := &http.Response{
		StatusCode: http.StatusBadGateway,
		Header: http.Header{
			"X-Request-Id": []string{"req-456"},
		},
	}

	err := adapters.NormalizeHTTPError(errors.New("upstream failed"), resp)

	var sandboxErr *sandboxerrors.SandboxError
	if !errors.As(err, &sandboxErr) {
		t.Fatal("expected normalized sandbox error")
	}
	if sandboxErr.RequestID != "req-456" {
		t.Fatalf("unexpected request id: %s", sandboxErr.RequestID)
	}
	if sandboxErr.StatusCode != http.StatusBadGateway {
		t.Fatalf("unexpected status code: %d", sandboxErr.StatusCode)
	}
}

func TestNormalizeHTTPErrorFallsBackToOriginalError(t *testing.T) {
	original := errors.New("original")
	err := adapters.NormalizeHTTPError(original, nil)
	if !errors.Is(err, original) {
		t.Fatal("expected original error to be preserved")
	}
}
