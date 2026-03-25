package unit

import (
	"context"
	"bytes"
	"net/http"
	"io"
	"testing"
	"time"

	"github.com/alibaba/opensandbox/sdks/sandbox/go/sandbox/config"
	"github.com/alibaba/opensandbox/sdks/sandbox/go/sandbox/factory"
	"github.com/alibaba/opensandbox/sdks/sandbox/go/sandbox/models"
)

func TestFactoryExposesLifecycleExecdAndEgressStacks(t *testing.T) {
	var _ factory.AdapterFactory = (*factory.DefaultAdapterFactory)(nil)
}

func TestDefaultFactoryFailsWithoutInjectedLifecycleClient(t *testing.T) {
	_, err := (&factory.DefaultAdapterFactory{}).CreateLifecycleStack(factory.CreateLifecycleStackOptions{})
	if err == nil {
		t.Fatal("expected missing lifecycle client error")
	}
}

func TestDefaultFactoryBuildsLifecycleStackFromConnectionConfig(t *testing.T) {
	var gotAPIKey string
	var gotCustomHeader string
	client := &http.Client{
		Transport: roundTripFunc(func(r *http.Request) (*http.Response, error) {
			gotAPIKey = r.Header.Get("OPEN-SANDBOX-API-KEY")
			gotCustomHeader = r.Header.Get("X-Custom-Header")
			body := []byte(`{"items":[],"pagination":{"page":1,"pageSize":10,"totalItems":0,"totalPages":0,"hasNextPage":false}}`)
			return &http.Response{
				StatusCode: http.StatusOK,
				Header: http.Header{
					"Content-Type": []string{"application/json"},
				},
				Body: io.NopCloser(bytes.NewReader(body)),
			}, nil
		}),
		Timeout: 5 * time.Second,
	}

	cfg := &config.ConnectionConfig{
		Domain:     "example.test",
		Protocol:   "http",
		APIKey:     "key-123",
		HTTPClient: client,
		Headers: map[string]string{
			"X-Custom-Header": "value-1",
		},
	}

	stack, err := (&factory.DefaultAdapterFactory{}).CreateLifecycleStack(factory.CreateLifecycleStackOptions{
		ConnectionConfig: cfg,
		LifecycleBaseURL: cfg.BaseURL(),
	})
	if err != nil {
		t.Fatalf("create lifecycle stack: %v", err)
	}

	if _, err := stack.Sandboxes.ListSandboxes(context.Background(), models.SandboxFilter{}); err != nil {
		t.Fatalf("list sandboxes: %v", err)
	}

	if gotAPIKey != "key-123" {
		t.Fatalf("unexpected api key header: %q", gotAPIKey)
	}
	if gotCustomHeader != "value-1" {
		t.Fatalf("unexpected custom header: %q", gotCustomHeader)
	}
}

type roundTripFunc func(*http.Request) (*http.Response, error)

func (f roundTripFunc) RoundTrip(req *http.Request) (*http.Response, error) {
	return f(req)
}
