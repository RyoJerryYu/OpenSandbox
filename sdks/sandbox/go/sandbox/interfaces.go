package sandbox

import (
	"context"
	"time"

	"github.com/alibaba/opensandbox/sdks/sandbox/go/sandbox/models"
	"github.com/alibaba/opensandbox/sdks/sandbox/go/sandbox/services"
)

// SandboxServices groups the service adapters exposed by a connected sandbox.
type SandboxServices struct {
	Commands services.ExecdCommands
	Files    services.SandboxFiles
	Health   services.ExecdHealth
	Metrics  services.ExecdMetrics
}

// Sandbox is the public interface for interacting with one sandbox instance.
type Sandbox interface {
	ID() string
	Services() SandboxServices
	GetInfo(ctx context.Context) (*models.SandboxInfo, error)
	GetEndpoint(ctx context.Context, port int) (*models.SandboxEndpoint, error)
	GetEndpointURL(ctx context.Context, port int) (string, error)
	IsHealthy(ctx context.Context) bool
	WaitUntilReady(ctx context.Context, opts *WaitUntilReadyOptions) error
	GetEgressPolicy(ctx context.Context) (*models.NetworkPolicy, error)
	PatchEgressRules(ctx context.Context, rules []models.NetworkRule) error
	Renew(ctx context.Context, timeout time.Duration) (*models.RenewSandboxExpirationResponse, error)
	Pause(ctx context.Context) error
	Kill(ctx context.Context) error
	Close() error
	Resume(ctx context.Context, opts ResumeOptions) (Sandbox, error)
}

// SandboxManager is the public interface for lifecycle-only sandbox management.
type SandboxManager interface {
	ListSandboxInfos(ctx context.Context, filter SandboxFilter) (*models.ListSandboxesResponse, error)
	GetSandboxInfo(ctx context.Context, sandboxID string) (*models.SandboxInfo, error)
	KillSandbox(ctx context.Context, sandboxID string) error
	PauseSandbox(ctx context.Context, sandboxID string) error
	ResumeSandbox(ctx context.Context, sandboxID string) error
	RenewSandbox(ctx context.Context, sandboxID string, timeout time.Duration) error
	Close() error
}
