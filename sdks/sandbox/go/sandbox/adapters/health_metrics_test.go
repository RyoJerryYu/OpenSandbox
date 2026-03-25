package adapters

import (
	"context"
	"net/http"
	"testing"

	execdapi "github.com/alibaba/opensandbox/sdks/sandbox/go/sandbox/internal/openapi/execd"
)

func TestHealthAdapterPingReturnsTrueOn200(t *testing.T) {
	adapter := NewHealthAdapter(&fakeExecdClient{
		pingResponse: &execdapi.PingResponse{
			HTTPResponse: &http.Response{StatusCode: http.StatusOK},
		},
	})

	ok, err := adapter.Ping(context.Background())
	if err != nil {
		t.Fatalf("ping: %v", err)
	}
	if !ok {
		t.Fatal("expected health ping to return true")
	}
}

func TestMetricsAdapterConvertsTimestampAndResourceNumbers(t *testing.T) {
	adapter := NewMetricsAdapter(&fakeExecdClient{
		metricsResponse: &execdapi.GetMetricsResponse{
			HTTPResponse: &http.Response{StatusCode: http.StatusOK},
			JSON200: &execdapi.Metrics{
				CpuCount:    4,
				CpuUsedPct:  12.5,
				MemTotalMib: 2048,
				MemUsedMib:  256,
				Timestamp:   1700000000000,
			},
		},
	})

	metrics, err := adapter.GetMetrics(context.Background())
	if err != nil {
		t.Fatalf("get metrics: %v", err)
	}
	if metrics.CPUCount != 4 {
		t.Fatalf("unexpected cpu count: %v", metrics.CPUCount)
	}
	if metrics.CPUUsedPercentage != 12.5 {
		t.Fatalf("unexpected cpu pct: %v", metrics.CPUUsedPercentage)
	}
	if metrics.MemoryTotalMiB != 2048 {
		t.Fatalf("unexpected mem total: %v", metrics.MemoryTotalMiB)
	}
	if metrics.MemoryUsedMiB != 256 {
		t.Fatalf("unexpected mem used: %v", metrics.MemoryUsedMiB)
	}
	if metrics.Timestamp != 1700000000000 {
		t.Fatalf("unexpected timestamp: %d", metrics.Timestamp)
	}
}

type fakeExecdClient struct {
	pingResponse    *execdapi.PingResponse
	metricsResponse *execdapi.GetMetricsResponse
}

func (f *fakeExecdClient) PingWithResponse(ctx context.Context, reqEditors ...execdapi.RequestEditorFn) (*execdapi.PingResponse, error) {
	return f.pingResponse, nil
}

func (f *fakeExecdClient) GetMetricsWithResponse(ctx context.Context, reqEditors ...execdapi.RequestEditorFn) (*execdapi.GetMetricsResponse, error) {
	return f.metricsResponse, nil
}
