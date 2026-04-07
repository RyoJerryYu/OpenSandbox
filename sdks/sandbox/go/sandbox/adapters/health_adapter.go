package adapters

import (
	"context"
	"fmt"

	execdapi "github.com/alibaba/opensandbox/sdks/sandbox/go/sandbox/internal/openapi/execd"
)

type HealthClient interface {
	PingWithResponse(ctx context.Context, reqEditors ...execdapi.RequestEditorFn) (*execdapi.PingResponse, error)
}

type HealthAdapter struct {
	client HealthClient
}

func NewHealthAdapter(client HealthClient) *HealthAdapter {
	return &HealthAdapter{client: client}
}

func (a *HealthAdapter) Ping(ctx context.Context) (bool, error) {
	resp, err := a.client.PingWithResponse(ctx)
	if err != nil {
		return false, err
	}
	if resp.StatusCode() >= 200 && resp.StatusCode() < 300 {
		return true, nil
	}
	return false, NormalizeHTTPError(fmt.Errorf("ping failed: unexpected status %d", resp.StatusCode()), resp.HTTPResponse)
}
