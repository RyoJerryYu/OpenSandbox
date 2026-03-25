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

func Create(ctx context.Context, opts SandboxCreateOptions) (*Sandbox, error) {
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
	return sbx, nil
}

func Connect(ctx context.Context, opts SandboxConnectOptions) (*Sandbox, error) {
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

	return connectConstructedSandbox(ctx, adapterFactory, connectionConfig, lifecycleStack.Sandboxes, opts.SandboxID)
}

func (s *Sandbox) Resume(ctx context.Context, opts ResumeOptions) (*Sandbox, error) {
	if err := s.sandboxes.ResumeSandbox(ctx, s.ID); err != nil {
		return nil, err
	}

	connectionConfig := opts.ConnectionConfig
	if connectionConfig == nil {
		connectionConfig = s.connectionConfig
	}
	return Connect(ctx, SandboxConnectOptions{
		ConnectionConfig: connectionConfig,
		AdapterFactory:   opts.AdapterFactory,
		SandboxID:        s.ID,
	})
}

func connectConstructedSandbox(ctx context.Context, adapterFactory factory.AdapterFactory, connectionConfig *config.ConnectionConfig, sandboxes services.Sandboxes, sandboxID string) (*Sandbox, error) {
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
	_, err = adapterFactory.CreateEgressStack(factory.CreateEgressStackOptions{
		ConnectionConfig: connectionConfig,
		EgressBaseURL:    endpointToBaseURL(connectionConfig, egressEndpoint),
	})
	if err != nil {
		return nil, err
	}

	return &Sandbox{
		ID:               sandboxID,
		Commands:         execdStack.Commands,
		Files:            execdStack.Files,
		Health:           execdStack.Health,
		Metrics:          execdStack.Metrics,
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
