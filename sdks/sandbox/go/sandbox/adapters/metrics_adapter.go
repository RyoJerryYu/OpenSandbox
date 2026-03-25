package adapters

import (
	"context"
	"fmt"

	"github.com/alibaba/opensandbox/sdks/sandbox/go/sandbox/internal/convert"
	execdapi "github.com/alibaba/opensandbox/sdks/sandbox/go/sandbox/internal/openapi/execd"
	"github.com/alibaba/opensandbox/sdks/sandbox/go/sandbox/models"
)

type MetricsClient interface {
	GetMetricsWithResponse(ctx context.Context, reqEditors ...execdapi.RequestEditorFn) (*execdapi.GetMetricsResponse, error)
}

type MetricsAdapter struct {
	client MetricsClient
}

func NewMetricsAdapter(client MetricsClient) *MetricsAdapter {
	return &MetricsAdapter{client: client}
}

func (a *MetricsAdapter) GetMetrics(ctx context.Context) (*models.SandboxMetrics, error) {
	resp, err := a.client.GetMetricsWithResponse(ctx)
	if err != nil {
		return nil, err
	}
	if resp.JSON200 == nil {
		return nil, NormalizeHTTPError(fmt.Errorf("get metrics failed: unexpected status %d", resp.StatusCode()), resp.HTTPResponse)
	}
	return convert.FromExecdMetrics(resp.JSON200), nil
}
