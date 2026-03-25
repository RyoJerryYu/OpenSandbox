package factory

import (
	"bytes"
	"context"
	"io"
	"net/http"
	"strings"
	"testing"
	"time"

	"github.com/alibaba/opensandbox/sdks/sandbox/go/sandbox/config"
)

func TestDefaultAdapterFactoryCreateExecdStackBuildsHealthAndMetricsAdapters(t *testing.T) {
	var requests []capturedRequest
	cfg := &config.ConnectionConfig{
		APIKey:  "test-api-key",
		Headers: map[string]string{"X-Test-Header": "test-value"},
		HTTPClient: &http.Client{
			Transport: roundTripFunc(func(req *http.Request) (*http.Response, error) {
				requests = append(requests, capturedRequest{
					Path:    req.URL.Path,
					APIKey:  req.Header.Get("OPEN-SANDBOX-API-KEY"),
					Header:  req.Header.Get("X-Test-Header"),
					Method:  req.Method,
					Timeout: req.Context().Err(),
				})

				switch req.URL.Path {
				case "/ping":
					return jsonResponse(http.StatusOK, ""), nil
				case "/metrics":
					return jsonResponse(http.StatusOK, `{"cpu_count":4,"cpu_used_pct":12.5,"mem_total_mib":2048,"mem_used_mib":256,"timestamp":1700000000000}`), nil
				default:
					return jsonResponse(http.StatusNotFound, `{"code":"not_found","message":"not found"}`), nil
				}
			}),
			Timeout: 2 * time.Second,
		},
	}

	stack, err := (&DefaultAdapterFactory{}).CreateExecdStack(CreateExecdStackOptions{
		ConnectionConfig: cfg,
		ExecdBaseURL:     "http://sandbox.example:44772",
	})
	if err != nil {
		t.Fatalf("create execd stack: %v", err)
	}
	if stack.Health == nil {
		t.Fatal("expected health adapter")
	}
	if stack.Metrics == nil {
		t.Fatal("expected metrics adapter")
	}
	if stack.Files == nil {
		t.Fatal("expected filesystem adapter")
	}
	if stack.Commands == nil {
		t.Fatal("expected commands adapter")
	}

	ok, err := stack.Health.Ping(context.Background())
	if err != nil {
		t.Fatalf("ping: %v", err)
	}
	if !ok {
		t.Fatal("expected ping to succeed")
	}

	metrics, err := stack.Metrics.GetMetrics(context.Background())
	if err != nil {
		t.Fatalf("metrics: %v", err)
	}
	if metrics.CPUCount != 4 {
		t.Fatalf("unexpected cpu count: %v", metrics.CPUCount)
	}

	if len(requests) != 2 {
		t.Fatalf("expected 2 requests, got %d", len(requests))
	}
	if requests[0].Path != "/ping" {
		t.Fatalf("unexpected first path: %s", requests[0].Path)
	}
	if requests[1].Path != "/metrics" {
		t.Fatalf("unexpected second path: %s", requests[1].Path)
	}
	for _, request := range requests {
		if request.APIKey != "test-api-key" {
			t.Fatalf("unexpected api key header: %q", request.APIKey)
		}
		if request.Header != "test-value" {
			t.Fatalf("unexpected custom header: %q", request.Header)
		}
		if request.Method != http.MethodGet {
			t.Fatalf("unexpected method: %s", request.Method)
		}
	}
}

func TestDefaultAdapterFactoryCreateEgressStackBuildsAdapter(t *testing.T) {
	client := &http.Client{
		Transport: roundTripFunc(func(r *http.Request) (*http.Response, error) {
			body := []byte(`{"mode":"deny_all","policy":{"defaultAction":"deny","egress":[{"action":"allow","target":"pypi.org"}]}}`)
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

	stack, err := (&DefaultAdapterFactory{}).CreateEgressStack(CreateEgressStackOptions{
		ConnectionConfig: &config.ConnectionConfig{
			Domain:     "example.test",
			Protocol:   "http",
			HTTPClient: client,
		},
		EgressBaseURL: "http://sandbox.example:18080",
	})
	if err != nil {
		t.Fatalf("create egress stack: %v", err)
	}
	if stack.Egress == nil {
		t.Fatal("expected egress adapter")
	}

	policy, err := stack.Egress.GetPolicy(context.Background())
	if err != nil {
		t.Fatalf("get policy: %v", err)
	}
	if policy.DefaultAction != "deny" {
		t.Fatalf("unexpected default action: %s", policy.DefaultAction)
	}
}

type capturedRequest struct {
	Path    string
	APIKey  string
	Header  string
	Method  string
	Timeout error
}

type roundTripFunc func(req *http.Request) (*http.Response, error)

func (f roundTripFunc) RoundTrip(req *http.Request) (*http.Response, error) {
	return f(req)
}

func jsonResponse(statusCode int, body string) *http.Response {
	if body == "" {
		body = "{}"
	}
	header := make(http.Header)
	header.Set("Content-Type", "application/json")
	return &http.Response{
		StatusCode: statusCode,
		Header:     header,
		Body:       io.NopCloser(strings.NewReader(body)),
	}
}
