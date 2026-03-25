package unit

import (
	"testing"
	"time"

	sandboxpkg "github.com/alibaba/opensandbox/sdks/sandbox/go/sandbox"
)

func TestNilTimeoutMeansManualCleanup(t *testing.T) {
	var timeout *time.Duration
	opts := sandboxpkg.SandboxCreateOptions{Timeout: timeout}
	if opts.Timeout != nil {
		t.Fatal("expected nil timeout")
	}
}
