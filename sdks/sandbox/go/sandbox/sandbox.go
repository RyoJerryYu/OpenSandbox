package sandbox

import (
	"context"
	"time"

	"github.com/alibaba/opensandbox/sdks/sandbox/go/sandbox/config"
	"github.com/alibaba/opensandbox/sdks/sandbox/go/sandbox/models"
	"github.com/alibaba/opensandbox/sdks/sandbox/go/sandbox/services"
)

type Sandbox struct {
	ID string

	Commands services.ExecdCommands
	Files    services.SandboxFiles
	Health   services.ExecdHealth
	Metrics  services.ExecdMetrics

	connectionConfig *config.ConnectionConfig
	sandboxes        services.Sandboxes
	closeFn          func() error
}

func (s *Sandbox) GetInfo(ctx context.Context) (*models.SandboxInfo, error) {
	return s.sandboxes.GetSandbox(ctx, s.ID)
}

func (s *Sandbox) GetEndpoint(ctx context.Context, port int) (*models.SandboxEndpoint, error) {
	return s.sandboxes.GetSandboxEndpoint(ctx, s.ID, port, s.connectionConfig.UseServerProxy)
}

func (s *Sandbox) Renew(ctx context.Context, timeout time.Duration) (*models.RenewSandboxExpirationResponse, error) {
	return s.sandboxes.RenewSandboxExpiration(ctx, s.ID, time.Now().UTC().Add(timeout))
}

func (s *Sandbox) Pause(ctx context.Context) error {
	return s.sandboxes.PauseSandbox(ctx, s.ID)
}

func (s *Sandbox) Kill(ctx context.Context) error {
	return s.sandboxes.DeleteSandbox(ctx, s.ID)
}

func (s *Sandbox) Close() error {
	if s == nil || s.closeFn == nil {
		return nil
	}
	return s.closeFn()
}
