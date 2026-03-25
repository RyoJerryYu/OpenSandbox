package transport

import "net/http"

type HTTPClientProvider struct {
	HTTPClient    *http.Client
	SSEHTTPClient *http.Client
}
