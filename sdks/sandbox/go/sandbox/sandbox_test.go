package sandbox

import (
	"context"
	"testing"
	"time"

	"github.com/alibaba/opensandbox/sdks/sandbox/go/sandbox/config"
	"github.com/alibaba/opensandbox/sdks/sandbox/go/sandbox/models"
)

func TestSandboxRenewUsesNowPlusTimeout(t *testing.T) {
	fake := &sandboxLifecycleFake{}
	sbx := &Sandbox{
		ID:               "sbx-1",
		connectionConfig: config.DefaultConnectionConfig(),
		sandboxes:        fake,
	}

	before := time.Now()
	if _, err := sbx.Renew(context.Background(), 90*time.Second); err != nil {
		t.Fatalf("renew sandbox: %v", err)
	}

	minExpected := before.Add(88 * time.Second)
	maxExpected := before.Add(92 * time.Second)
	if fake.renewExpiresAt.Before(minExpected) || fake.renewExpiresAt.After(maxExpected) {
		t.Fatalf("unexpected renew time: %s", fake.renewExpiresAt)
	}
}

func TestSandboxGetEndpointUsesConnectionProxyFlag(t *testing.T) {
	fake := &sandboxLifecycleFake{
		endpointResponse: &models.SandboxEndpoint{Endpoint: "domain/proxy/44772"},
	}
	sbx := &Sandbox{
		ID: "sbx-1",
		connectionConfig: &config.ConnectionConfig{
			Domain:         config.DefaultDomain,
			Protocol:       config.DefaultProtocol,
			UseServerProxy: true,
		},
		sandboxes: fake,
	}

	if _, err := sbx.GetEndpoint(context.Background(), 44772); err != nil {
		t.Fatalf("get endpoint: %v", err)
	}
	if !fake.useServerProxy {
		t.Fatal("expected useServerProxy to be forwarded")
	}
}

func TestSandboxCloseCallsCloseFunc(t *testing.T) {
	closed := false
	sbx := &Sandbox{
		closeFn: func() error {
			closed = true
			return nil
		},
	}

	if err := sbx.Close(); err != nil {
		t.Fatalf("close sandbox: %v", err)
	}
	if !closed {
		t.Fatal("expected close function to be called")
	}
}

type sandboxLifecycleFake struct {
	endpointResponse *models.SandboxEndpoint
	renewExpiresAt   time.Time
	useServerProxy   bool
}

func (f *sandboxLifecycleFake) CreateSandbox(context.Context, models.CreateSandboxRequest) (*models.CreateSandboxResponse, error) {
	return nil, nil
}

func (f *sandboxLifecycleFake) GetSandbox(context.Context, string) (*models.SandboxInfo, error) {
	return nil, nil
}

func (f *sandboxLifecycleFake) ListSandboxes(context.Context, models.SandboxFilter) (*models.ListSandboxesResponse, error) {
	return nil, nil
}

func (f *sandboxLifecycleFake) DeleteSandbox(context.Context, string) error {
	return nil
}

func (f *sandboxLifecycleFake) PauseSandbox(context.Context, string) error {
	return nil
}

func (f *sandboxLifecycleFake) ResumeSandbox(context.Context, string) error {
	return nil
}

func (f *sandboxLifecycleFake) RenewSandboxExpiration(_ context.Context, _ string, expiresAt time.Time) (*models.RenewSandboxExpirationResponse, error) {
	f.renewExpiresAt = expiresAt
	return &models.RenewSandboxExpirationResponse{ExpiresAt: &expiresAt}, nil
}

func (f *sandboxLifecycleFake) GetSandboxEndpoint(_ context.Context, _ string, _ int, useServerProxy bool) (*models.SandboxEndpoint, error) {
	f.useServerProxy = useServerProxy
	return f.endpointResponse, nil
}
