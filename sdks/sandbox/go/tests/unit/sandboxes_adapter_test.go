package unit

import (
	"context"
	"testing"
	"time"

	"github.com/alibaba/opensandbox/sdks/sandbox/go/sandbox/adapters"
	"github.com/alibaba/opensandbox/sdks/sandbox/go/sandbox/models"
)

func TestSandboxesAdapterConvertsCreateResponse(t *testing.T) {
	now := time.Now().UTC()
	client := &fakeLifecycleClient{
		createResponse: &models.CreateSandboxResponse{
			ID: "sbx-123",
			Status: models.SandboxStatus{
				State: "Running",
			},
			CreatedAt: now,
			ExpiresAt: &now,
		},
	}

	adapter := adapters.NewSandboxesAdapter(client)
	resp, err := adapter.CreateSandbox(context.Background(), models.CreateSandboxRequest{})
	if err != nil {
		t.Fatalf("create sandbox: %v", err)
	}
	if resp.ID != "sbx-123" {
		t.Fatalf("unexpected id: %s", resp.ID)
	}
	if resp.Status.State != "Running" {
		t.Fatalf("unexpected state: %s", resp.Status.State)
	}
}

func TestSandboxesAdapterPassesUseServerProxyToEndpointLookup(t *testing.T) {
	client := &fakeLifecycleClient{
		endpointResponse: &models.SandboxEndpoint{
			Endpoint: "localhost/proxy/44772",
		},
	}

	adapter := adapters.NewSandboxesAdapter(client)
	_, err := adapter.GetSandboxEndpoint(context.Background(), "sbx-123", 44772, true)
	if err != nil {
		t.Fatalf("get endpoint: %v", err)
	}

	if !client.useServerProxy {
		t.Fatal("expected useServerProxy to be forwarded")
	}
	if client.endpointSandboxID != "sbx-123" {
		t.Fatalf("unexpected sandbox id: %s", client.endpointSandboxID)
	}
	if client.endpointPort != 44772 {
		t.Fatalf("unexpected port: %d", client.endpointPort)
	}
}

type fakeLifecycleClient struct {
	createResponse     *models.CreateSandboxResponse
	endpointResponse   *models.SandboxEndpoint
	endpointSandboxID  string
	endpointPort       int
	useServerProxy     bool
}

func (f *fakeLifecycleClient) CreateSandbox(context.Context, models.CreateSandboxRequest) (*models.CreateSandboxResponse, error) {
	return f.createResponse, nil
}

func (f *fakeLifecycleClient) GetSandbox(context.Context, string) (*models.SandboxInfo, error) {
	return nil, nil
}

func (f *fakeLifecycleClient) ListSandboxes(context.Context, models.SandboxFilter) (*models.ListSandboxesResponse, error) {
	return nil, nil
}

func (f *fakeLifecycleClient) DeleteSandbox(context.Context, string) error {
	return nil
}

func (f *fakeLifecycleClient) PauseSandbox(context.Context, string) error {
	return nil
}

func (f *fakeLifecycleClient) ResumeSandbox(context.Context, string) error {
	return nil
}

func (f *fakeLifecycleClient) RenewSandboxExpiration(context.Context, string, time.Time) (*models.RenewSandboxExpirationResponse, error) {
	return nil, nil
}

func (f *fakeLifecycleClient) GetSandboxEndpoint(_ context.Context, sandboxID string, port int, useServerProxy bool) (*models.SandboxEndpoint, error) {
	f.endpointSandboxID = sandboxID
	f.endpointPort = port
	f.useServerProxy = useServerProxy
	return f.endpointResponse, nil
}
