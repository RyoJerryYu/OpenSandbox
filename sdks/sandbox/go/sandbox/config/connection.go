package config

import (
	"net/http"
	"time"

	"github.com/alibaba/opensandbox/sdks/sandbox/go/sandbox/internal/transport"
)

type ConnectionConfig struct {
	// Domain is the lifecycle API host or host:port without scheme.
	Domain         string
	// Protocol is the scheme used for lifecycle and sandbox endpoint requests.
	Protocol       string
	// APIKey is forwarded to lifecycle requests when the deployment requires API-key auth.
	APIKey         string
	// Headers are merged into every outbound SDK request.
	Headers        map[string]string
	// RequestTimeout bounds non-streaming HTTP calls made by SDK-managed clients.
	RequestTimeout time.Duration
	// UseServerProxy asks the lifecycle API to return server-proxied sandbox endpoints.
	UseServerProxy bool
	// HTTPClient is used for regular request/response traffic. Caller-owned when provided explicitly.
	HTTPClient     *http.Client
	// SSEHTTPClient is used for streaming requests such as SSE command execution.
	SSEHTTPClient  *http.Client

	ownsHTTPClient    bool
	ownsSSEHTTPClient bool
}

// DefaultConnectionConfig returns a local-development oriented connection configuration.
func DefaultConnectionConfig() *ConnectionConfig {
	return &ConnectionConfig{
		Domain:         DefaultDomain,
		Protocol:       DefaultProtocol,
		Headers:        map[string]string{},
		RequestTimeout: 30 * time.Second,
	}
}

// BaseURL returns the lifecycle API base URL derived from Protocol and Domain.
func (c *ConnectionConfig) BaseURL() string {
	return transport.BaseURL(c.Protocol, c.Domain)
}

// CloneWithHTTPClientIfMissing returns a copy with default HTTP clients when none are configured.
func (c *ConnectionConfig) CloneWithHTTPClientIfMissing() *ConnectionConfig {
	if c == nil {
		c = DefaultConnectionConfig()
	}

	clone := *c
	if clone.Headers == nil {
		clone.Headers = map[string]string{}
	}
	if clone.HTTPClient == nil {
		clone.HTTPClient = &http.Client{
			Timeout: clone.RequestTimeout,
		}
		clone.ownsHTTPClient = true
	}
	if clone.SSEHTTPClient == nil {
		clone.SSEHTTPClient = &http.Client{
			Timeout: clone.RequestTimeout,
		}
		clone.ownsSSEHTTPClient = true
	}
	return &clone
}

// Close releases idle connections owned by SDK-managed HTTP clients.
func (c *ConnectionConfig) Close() error {
	if c == nil {
		return nil
	}

	if c.ownsHTTPClient {
		closeHTTPClient(c.HTTPClient)
	}
	if c.ownsSSEHTTPClient && c.SSEHTTPClient != c.HTTPClient {
		closeHTTPClient(c.SSEHTTPClient)
	}
	return nil
}

type idleCloser interface {
	CloseIdleConnections()
}

func closeHTTPClient(client *http.Client) {
	if client == nil || client.Transport == nil {
		return
	}
	if closer, ok := client.Transport.(idleCloser); ok {
		closer.CloseIdleConnections()
	}
}
