package sandbox

import (
	"context"
	"time"

	"github.com/alibaba/opensandbox/sdks/sandbox/go/sandbox/config"
	"github.com/alibaba/opensandbox/sdks/sandbox/go/sandbox/factory"
	"github.com/alibaba/opensandbox/sdks/sandbox/go/sandbox/models"
	"github.com/alibaba/opensandbox/sdks/sandbox/go/sandbox/services"
)

// SandboxManagerOptions configures construction of a SandboxManager.
type SandboxManagerOptions struct {
	ConnectionConfig *config.ConnectionConfig
	AdapterFactory   factory.AdapterFactory
	CloseFunc        func() error
}

// SandboxManager provides lifecycle operations without attaching execd adapters.
type SandboxManager struct {
	sandboxes services.Sandboxes
	closeFn   func() error
}

// NewSandboxManager constructs a lifecycle-only manager using the provided connection settings.
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

// ListSandboxInfos returns sandboxes visible to the current connection configuration.
func (m *SandboxManager) ListSandboxInfos(ctx context.Context, filter SandboxFilter) (*models.ListSandboxesResponse, error) {
	return m.sandboxes.ListSandboxes(ctx, filter)
}

// GetSandboxInfo fetches lifecycle information for one sandbox.
func (m *SandboxManager) GetSandboxInfo(ctx context.Context, sandboxID string) (*models.SandboxInfo, error) {
	return m.sandboxes.GetSandbox(ctx, sandboxID)
}

// KillSandbox deletes a sandbox remotely.
func (m *SandboxManager) KillSandbox(ctx context.Context, sandboxID string) error {
	return m.sandboxes.DeleteSandbox(ctx, sandboxID)
}

// PauseSandbox pauses a sandbox remotely.
func (m *SandboxManager) PauseSandbox(ctx context.Context, sandboxID string) error {
	return m.sandboxes.PauseSandbox(ctx, sandboxID)
}

// ResumeSandbox resumes a paused sandbox remotely.
func (m *SandboxManager) ResumeSandbox(ctx context.Context, sandboxID string) error {
	return m.sandboxes.ResumeSandbox(ctx, sandboxID)
}

// RenewSandbox extends a sandbox expiration relative to the current time.
func (m *SandboxManager) RenewSandbox(ctx context.Context, sandboxID string, timeout time.Duration) error {
	_, err := m.sandboxes.RenewSandboxExpiration(ctx, sandboxID, time.Now().UTC().Add(timeout))
	return err
}

// Close releases local resources owned by the manager.
func (m *SandboxManager) Close() error {
	if m == nil || m.closeFn == nil {
		return nil
	}
	return m.closeFn()
}
