package transport_test

import (
	"net/http"
	"testing"

	"github.com/alibaba/opensandbox/sdks/sandbox/go/sandbox/config"
	"github.com/alibaba/opensandbox/sdks/sandbox/go/sandbox/internal/transport"
)

func TestMergeHeadersPrefersEndpointHeaders(t *testing.T) {
	merged := transport.MergeHeaders(
		map[string]string{"X-Test": "base", "X-Base": "keep"},
		map[string]string{"X-Test": "endpoint", "X-Endpoint": "keep"},
	)

	if got := merged["X-Test"]; got != "endpoint" {
		t.Fatalf("expected endpoint override, got %q", got)
	}
	if got := merged["X-Base"]; got != "keep" {
		t.Fatalf("expected base header to survive, got %q", got)
	}
	if got := merged["X-Endpoint"]; got != "keep" {
		t.Fatalf("expected endpoint header to exist, got %q", got)
	}
}

func TestCloseOnlyClosesOwnedClients(t *testing.T) {
	userTransport := &closingTransport{}
	cfg := &config.ConnectionConfig{
		Domain:     config.DefaultDomain,
		HTTPClient: &http.Client{Transport: userTransport},
	}

	cloned := cfg.CloneWithHTTPClientIfMissing()
	if err := cloned.Close(); err != nil {
		t.Fatalf("close config: %v", err)
	}
	if userTransport.closed {
		t.Fatal("expected user supplied transport to stay open")
	}

	owned := config.DefaultConnectionConfig().CloneWithHTTPClientIfMissing()
	if err := owned.Close(); err != nil {
		t.Fatalf("close owned config: %v", err)
	}
}

type closingTransport struct {
	closed bool
}

func (t *closingTransport) RoundTrip(*http.Request) (*http.Response, error) {
	return nil, nil
}

func (t *closingTransport) CloseIdleConnections() {
	t.closed = true
}
