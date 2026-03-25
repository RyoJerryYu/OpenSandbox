package sandbox

import (
	"context"
	"time"

	"github.com/alibaba/opensandbox/sdks/sandbox/go/sandbox/config"
	"github.com/alibaba/opensandbox/sdks/sandbox/go/sandbox/factory"
	"github.com/alibaba/opensandbox/sdks/sandbox/go/sandbox/models"
	"github.com/alibaba/opensandbox/sdks/sandbox/go/sandbox/services"
)

type SandboxManagerOptions struct {
	ConnectionConfig *config.ConnectionConfig
	AdapterFactory   factory.AdapterFactory
	CloseFunc        func() error
}

type SandboxManager struct {
	sandboxes services.Sandboxes
	closeFn   func() error
}

func NewSandboxManager(opts SandboxManagerOptions) (*SandboxManager, error) {
	connectionConfig := opts.ConnectionConfig
	if connectionConfig == nil {
		connectionConfig = config.DefaultConnectionConfig()
	}
	connectionConfig = connectionConfig.CloneWithHTTPClientIfMissing()

	adapterFactory := opts.AdapterFactory
	if adapterFactory == nil {
		adapterFactory = &factory.DefaultAdapterFactory{}
	}

	stack, err := adapterFactory.CreateLifecycleStack(factory.CreateLifecycleStackOptions{
		ConnectionConfig: connectionConfig,
		LifecycleBaseURL: connectionConfig.BaseURL(),
	})
	if err != nil {
		_ = connectionConfig.Close()
		return nil, err
	}

	closeFn := opts.CloseFunc
	if closeFn == nil {
		closeFn = connectionConfig.Close
	}

	return &SandboxManager{
		sandboxes: stack.Sandboxes,
		closeFn:   closeFn,
	}, nil
}

func (m *SandboxManager) ListSandboxInfos(ctx context.Context, filter SandboxFilter) (*models.ListSandboxesResponse, error) {
	return m.sandboxes.ListSandboxes(ctx, filter)
}

func (m *SandboxManager) GetSandboxInfo(ctx context.Context, sandboxID string) (*models.SandboxInfo, error) {
	return m.sandboxes.GetSandbox(ctx, sandboxID)
}

func (m *SandboxManager) KillSandbox(ctx context.Context, sandboxID string) error {
	return m.sandboxes.DeleteSandbox(ctx, sandboxID)
}

func (m *SandboxManager) PauseSandbox(ctx context.Context, sandboxID string) error {
	return m.sandboxes.PauseSandbox(ctx, sandboxID)
}

func (m *SandboxManager) ResumeSandbox(ctx context.Context, sandboxID string) error {
	return m.sandboxes.ResumeSandbox(ctx, sandboxID)
}

func (m *SandboxManager) RenewSandbox(ctx context.Context, sandboxID string, timeout time.Duration) error {
	_, err := m.sandboxes.RenewSandboxExpiration(ctx, sandboxID, time.Now().UTC().Add(timeout))
	return err
}

func (m *SandboxManager) Close() error {
	if m == nil || m.closeFn == nil {
		return nil
	}
	return m.closeFn()
}
