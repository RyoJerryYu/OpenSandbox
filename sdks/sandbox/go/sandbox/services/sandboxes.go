package services

import (
	"context"
	"time"

	"github.com/alibaba/opensandbox/sdks/sandbox/go/sandbox/models"
)

// Sandboxes defines lifecycle operations exposed by the OpenSandbox server.
type Sandboxes interface {
	CreateSandbox(ctx context.Context, req models.CreateSandboxRequest) (*models.CreateSandboxResponse, error)
	GetSandbox(ctx context.Context, sandboxID string) (*models.SandboxInfo, error)
	ListSandboxes(ctx context.Context, filter models.SandboxFilter) (*models.ListSandboxesResponse, error)
	DeleteSandbox(ctx context.Context, sandboxID string) error
	PauseSandbox(ctx context.Context, sandboxID string) error
	ResumeSandbox(ctx context.Context, sandboxID string) error
	RenewSandboxExpiration(ctx context.Context, sandboxID string, expiresAt time.Time) (*models.RenewSandboxExpirationResponse, error)
	GetSandboxEndpoint(ctx context.Context, sandboxID string, port int, useServerProxy bool) (*models.SandboxEndpoint, error)
}
