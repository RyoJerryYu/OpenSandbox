package adapters

import (
	"context"
	"time"

	"github.com/alibaba/opensandbox/sdks/sandbox/go/sandbox/internal/convert"
	"github.com/alibaba/opensandbox/sdks/sandbox/go/sandbox/models"
)

type LifecycleClient interface {
	CreateSandbox(ctx context.Context, req models.CreateSandboxRequest) (*models.CreateSandboxResponse, error)
	GetSandbox(ctx context.Context, sandboxID string) (*models.SandboxInfo, error)
	ListSandboxes(ctx context.Context, filter models.SandboxFilter) (*models.ListSandboxesResponse, error)
	DeleteSandbox(ctx context.Context, sandboxID string) error
	PauseSandbox(ctx context.Context, sandboxID string) error
	ResumeSandbox(ctx context.Context, sandboxID string) error
	RenewSandboxExpiration(ctx context.Context, sandboxID string, expiresAt time.Time) (*models.RenewSandboxExpirationResponse, error)
	GetSandboxEndpoint(ctx context.Context, sandboxID string, port int, useServerProxy bool) (*models.SandboxEndpoint, error)
}

type SandboxesAdapter struct {
	client LifecycleClient
}

func NewSandboxesAdapter(client LifecycleClient) *SandboxesAdapter {
	return &SandboxesAdapter{client: client}
}

func (a *SandboxesAdapter) CreateSandbox(ctx context.Context, req models.CreateSandboxRequest) (*models.CreateSandboxResponse, error) {
	resp, err := a.client.CreateSandbox(ctx, req)
	if err != nil {
		return nil, err
	}
	return convert.CloneCreateSandboxResponse(resp), nil
}

func (a *SandboxesAdapter) GetSandbox(ctx context.Context, sandboxID string) (*models.SandboxInfo, error) {
	resp, err := a.client.GetSandbox(ctx, sandboxID)
	if err != nil {
		return nil, err
	}
	return convert.CloneSandboxInfo(resp), nil
}

func (a *SandboxesAdapter) ListSandboxes(ctx context.Context, filter models.SandboxFilter) (*models.ListSandboxesResponse, error) {
	resp, err := a.client.ListSandboxes(ctx, filter)
	if err != nil {
		return nil, err
	}
	return convert.CloneListSandboxesResponse(resp), nil
}

func (a *SandboxesAdapter) DeleteSandbox(ctx context.Context, sandboxID string) error {
	return a.client.DeleteSandbox(ctx, sandboxID)
}

func (a *SandboxesAdapter) PauseSandbox(ctx context.Context, sandboxID string) error {
	return a.client.PauseSandbox(ctx, sandboxID)
}

func (a *SandboxesAdapter) ResumeSandbox(ctx context.Context, sandboxID string) error {
	return a.client.ResumeSandbox(ctx, sandboxID)
}

func (a *SandboxesAdapter) RenewSandboxExpiration(ctx context.Context, sandboxID string, expiresAt time.Time) (*models.RenewSandboxExpirationResponse, error) {
	resp, err := a.client.RenewSandboxExpiration(ctx, sandboxID, expiresAt)
	if err != nil {
		return nil, err
	}
	return convert.CloneRenewSandboxExpirationResponse(resp), nil
}

func (a *SandboxesAdapter) GetSandboxEndpoint(ctx context.Context, sandboxID string, port int, useServerProxy bool) (*models.SandboxEndpoint, error) {
	resp, err := a.client.GetSandboxEndpoint(ctx, sandboxID, port, useServerProxy)
	if err != nil {
		return nil, err
	}
	return convert.CloneSandboxEndpoint(resp), nil
}
