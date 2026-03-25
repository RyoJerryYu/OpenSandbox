package config

import (
	"net/http"
	"time"

	"github.com/alibaba/opensandbox/sdks/sandbox/go/sandbox/internal/transport"
)

type ConnectionConfig struct {
	Domain         string
	Protocol       string
	APIKey         string
	Headers        map[string]string
	RequestTimeout time.Duration
	UseServerProxy bool
	HTTPClient     *http.Client
	SSEHTTPClient  *http.Client

	ownsHTTPClient    bool
	ownsSSEHTTPClient bool
}

func DefaultConnectionConfig() *ConnectionConfig {
	return &ConnectionConfig{
		Domain:         DefaultDomain,
		Protocol:       DefaultProtocol,
		Headers:        map[string]string{},
		RequestTimeout: 30 * time.Second,
	}
}

func (c *ConnectionConfig) BaseURL() string {
	return transport.BaseURL(c.Protocol, c.Domain)
}

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
