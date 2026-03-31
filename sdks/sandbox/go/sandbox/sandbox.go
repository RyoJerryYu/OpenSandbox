package sandbox

import (
	"context"
	"fmt"
	"time"

	"github.com/alibaba/opensandbox/sdks/sandbox/go/sandbox/config"
	"github.com/alibaba/opensandbox/sdks/sandbox/go/sandbox/factory"
	"github.com/alibaba/opensandbox/sdks/sandbox/go/sandbox/models"
	"github.com/alibaba/opensandbox/sdks/sandbox/go/sandbox/services"
)

const (
	defaultExecdPort  = 44772
	defaultEgressPort = 18080
)

// SandboxImpl is the concrete implementation behind the public Sandbox interface.
type SandboxImpl struct {
	id               string
	services         SandboxServices
	egress           services.Egress
	connectionConfig *config.ConnectionConfig
	sandboxes        services.Sandboxes
	closeFn          func() error
}

// Create provisions a new sandbox and, by default, waits for it to become ready.
func Create(ctx context.Context, opts SandboxCreateOptions) (Sandbox, error) {
	connectionConfig := opts.ConnectionConfig
	if connectionConfig == nil {
		connectionConfig = config.DefaultConnectionConfig()
	}
	connectionConfig = connectionConfig.CloneWithHTTPClientIfMissing()

	adapterFactory := opts.AdapterFactory
	if adapterFactory == nil {
		adapterFactory = &factory.DefaultAdapterFactory{}
	}

	lifecycleStack, err := adapterFactory.CreateLifecycleStack(factory.CreateLifecycleStackOptions{
		ConnectionConfig: connectionConfig,
		LifecycleBaseURL: connectionConfig.BaseURL(),
	})
	if err != nil {
		_ = connectionConfig.Close()
		return nil, err
	}

	var timeoutSeconds *int
	if opts.Timeout != nil {
		v := int(opts.Timeout.Seconds())
		timeoutSeconds = &v
	}

	created, err := lifecycleStack.Sandboxes.CreateSandbox(ctx, models.CreateSandboxRequest{
		Image:          opts.Image,
		Entrypoint:     []string{"tail", "-f", "/dev/null"},
		Timeout:        timeoutSeconds,
		ResourceLimits: map[string]string{},
	})
	if err != nil {
		_ = connectionConfig.Close()
		return nil, err
	}

	sbx, err := connectConstructedSandbox(ctx, adapterFactory, connectionConfig, lifecycleStack.Sandboxes, created.ID)
	if err != nil {
		_ = lifecycleStack.Sandboxes.DeleteSandbox(ctx, created.ID)
		_ = connectionConfig.Close()
		return nil, err
	}
	if err := waitIfNeeded(ctx, sbx, readinessOptionsFromCreate(opts)); err != nil {
		_ = lifecycleStack.Sandboxes.DeleteSandbox(ctx, created.ID)
		_ = connectionConfig.Close()
		return nil, err
	}
	return sbx, nil
}

// Connect attaches to an existing sandbox and, by default, waits for it to become ready.
func Connect(ctx context.Context, opts SandboxConnectOptions) (Sandbox, error) {
	connectionConfig := opts.ConnectionConfig
	if connectionConfig == nil {
		connectionConfig = config.DefaultConnectionConfig()
	}
	connectionConfig = connectionConfig.CloneWithHTTPClientIfMissing()

	adapterFactory := opts.AdapterFactory
	if adapterFactory == nil {
		adapterFactory = &factory.DefaultAdapterFactory{}
	}

	lifecycleStack, err := adapterFactory.CreateLifecycleStack(factory.CreateLifecycleStackOptions{
		ConnectionConfig: connectionConfig,
		LifecycleBaseURL: connectionConfig.BaseURL(),
	})
	if err != nil {
		_ = connectionConfig.Close()
		return nil, err
	}

	sbx, err := connectConstructedSandbox(ctx, adapterFactory, connectionConfig, lifecycleStack.Sandboxes, opts.SandboxID)
	if err != nil {
		_ = connectionConfig.Close()
		return nil, err
	}
	if err := waitIfNeeded(ctx, sbx, readinessOptionsFromConnect(opts)); err != nil {
		_ = connectionConfig.Close()
		return nil, err
	}
	return sbx, nil
}

// Resume resumes the current sandbox remotely and returns a fresh connected Sandbox handle.
func (s *SandboxImpl) Resume(ctx context.Context, opts ResumeOptions) (Sandbox, error) {
	if err := s.sandboxes.ResumeSandbox(ctx, s.id); err != nil {
		return nil, err
	}

	connectionConfig := opts.ConnectionConfig
	if connectionConfig == nil {
		connectionConfig = s.connectionConfig
	}
	return Connect(ctx, SandboxConnectOptions{
		ConnectionConfig: connectionConfig,
		AdapterFactory:   opts.AdapterFactory,
		SandboxID:        s.id,
		SkipHealthCheck:  opts.SkipHealthCheck,
		ReadyTimeout:     opts.ReadyTimeout,
		PollInterval:     opts.PollInterval,
		HealthCheck:      opts.HealthCheck,
	})
}

func connectConstructedSandbox(ctx context.Context, adapterFactory factory.AdapterFactory, connectionConfig *config.ConnectionConfig, sandboxes services.Sandboxes, sandboxID string) (*SandboxImpl, error) {
	execdEndpoint, err := sandboxes.GetSandboxEndpoint(ctx, sandboxID, defaultExecdPort, connectionConfig.UseServerProxy)
	if err != nil {
		return nil, err
	}
	egressEndpoint, err := sandboxes.GetSandboxEndpoint(ctx, sandboxID, defaultEgressPort, connectionConfig.UseServerProxy)
	if err != nil {
		return nil, err
	}

	execdStack, err := adapterFactory.CreateExecdStack(factory.CreateExecdStackOptions{
		ConnectionConfig: connectionConfig,
		ExecdBaseURL:     endpointToBaseURL(connectionConfig, execdEndpoint),
	})
	if err != nil {
		return nil, err
	}
	egressStack, err := adapterFactory.CreateEgressStack(factory.CreateEgressStackOptions{
		ConnectionConfig: connectionConfig,
		EgressBaseURL:    endpointToBaseURL(connectionConfig, egressEndpoint),
	})
	if err != nil {
		return nil, err
	}

	return &SandboxImpl{
		id: sandboxID,
		services: SandboxServices{
			Commands: execdStack.Commands,
			Files:    execdStack.Files,
			Health:   execdStack.Health,
			Metrics:  execdStack.Metrics,
		},
		egress:           egressStack.Egress,
		connectionConfig: connectionConfig,
		sandboxes:        sandboxes,
		closeFn:          connectionConfig.Close,
	}, nil
}

func endpointToBaseURL(connectionConfig *config.ConnectionConfig, endpoint *models.SandboxEndpoint) string {
	if endpoint == nil {
		return ""
	}
	return fmt.Sprintf("%s://%s", connectionConfig.Protocol, endpoint.Endpoint)
}

// GetInfo fetches the latest lifecycle information for the sandbox.
func (s *SandboxImpl) ID() string {
	if s == nil {
		return ""
	}
	return s.id
}

func (s *SandboxImpl) Services() SandboxServices {
	if s == nil {
		return SandboxServices{}
	}
	return s.services
}

// GetInfo fetches the latest lifecycle information for the sandbox.
func (s *SandboxImpl) GetInfo(ctx context.Context) (*models.SandboxInfo, error) {
	return s.sandboxes.GetSandbox(ctx, s.id)
}

// GetEndpoint resolves the endpoint for a service listening on the given port.
func (s *SandboxImpl) GetEndpoint(ctx context.Context, port int) (*models.SandboxEndpoint, error) {
	return s.sandboxes.GetSandboxEndpoint(ctx, s.id, port, s.connectionConfig.UseServerProxy)
}

// GetEndpointURL resolves the endpoint for a port and prefixes it with the configured scheme.
func (s *SandboxImpl) GetEndpointURL(ctx context.Context, port int) (string, error) {
	endpoint, err := s.GetEndpoint(ctx, port)
	if err != nil {
		return "", err
	}
	return endpointToBaseURL(s.connectionConfig, endpoint), nil
}

// IsHealthy reports whether the sandbox execd health check succeeds.
func (s *SandboxImpl) IsHealthy(ctx context.Context) bool {
	if s == nil || s.services.Health == nil {
		return false
	}
	ok, err := s.services.Health.Ping(ctx)
	if err != nil {
		return false
	}
	return ok
}

// GetEgressPolicy fetches the current egress policy when the egress sidecar is available.
func (s *SandboxImpl) GetEgressPolicy(ctx context.Context) (*models.NetworkPolicy, error) {
	if s == nil || s.egress == nil {
		return nil, nil
	}
	return s.egress.GetPolicy(ctx)
}

// PatchEgressRules patches egress rules using sidecar merge semantics.
func (s *SandboxImpl) PatchEgressRules(ctx context.Context, rules []models.NetworkRule) error {
	if s == nil || s.egress == nil {
		return nil
	}
	return s.egress.PatchRules(ctx, rules)
}

// Renew extends the sandbox expiration relative to the current time.
func (s *SandboxImpl) Renew(ctx context.Context, timeout time.Duration) (*models.RenewSandboxExpirationResponse, error) {
	return s.sandboxes.RenewSandboxExpiration(ctx, s.id, time.Now().UTC().Add(timeout))
}

// Pause pauses the sandbox remotely.
func (s *SandboxImpl) Pause(ctx context.Context) error {
	return s.sandboxes.PauseSandbox(ctx, s.id)
}

// Kill deletes the sandbox remotely.
func (s *SandboxImpl) Kill(ctx context.Context) error {
	return s.sandboxes.DeleteSandbox(ctx, s.id)
}

// Close releases local resources owned by this Sandbox handle.
func (s *SandboxImpl) Close() error {
	if s == nil || s.closeFn == nil {
		return nil
	}
	return s.closeFn()
}

func waitIfNeeded(ctx context.Context, sbx *SandboxImpl, opts *WaitUntilReadyOptions) error {
	if sbx == nil || opts == nil {
		return nil
	}
	if opts.SkipStatePoll && opts.CustomHealthCheck == nil && sbx.services.Health == nil {
		return nil
	}
	return sbx.WaitUntilReady(ctx, opts)
}

func readinessOptionsFromCreate(opts SandboxCreateOptions) *WaitUntilReadyOptions {
	if opts.SkipHealthCheck {
		return nil
	}
	return &WaitUntilReadyOptions{
		Timeout:           opts.ReadyTimeout,
		PollInterval:      opts.PollInterval,
		CustomHealthCheck: opts.HealthCheck,
	}
}

func readinessOptionsFromConnect(opts SandboxConnectOptions) *WaitUntilReadyOptions {
	if opts.SkipHealthCheck {
		return nil
	}
	return &WaitUntilReadyOptions{
		Timeout:           opts.ReadyTimeout,
		PollInterval:      opts.PollInterval,
		CustomHealthCheck: opts.HealthCheck,
	}
}
