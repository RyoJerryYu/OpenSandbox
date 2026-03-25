package unit

import (
	"testing"

	"github.com/alibaba/opensandbox/sdks/sandbox/go/sandbox/config"
)

func TestDefaultConnectionConfigUsesLocalhost(t *testing.T) {
	cfg := config.DefaultConnectionConfig()
	if got := cfg.BaseURL(); got != "http://localhost:8080" {
		t.Fatalf("unexpected base url: %s", got)
	}
}
