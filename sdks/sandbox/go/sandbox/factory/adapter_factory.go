package factory

import (
	"github.com/alibaba/opensandbox/sdks/sandbox/go/sandbox/adapters"
	"github.com/alibaba/opensandbox/sdks/sandbox/go/sandbox/config"
	"github.com/alibaba/opensandbox/sdks/sandbox/go/sandbox/services"
)

type LifecycleStack struct {
	Sandboxes services.Sandboxes
}

type ExecdStack struct {
	Commands services.ExecdCommands
	Files    services.SandboxFiles
	Health   services.ExecdHealth
	Metrics  services.ExecdMetrics
}

type EgressStack struct {
	Egress services.Egress
}

type CreateLifecycleStackOptions struct {
	ConnectionConfig *config.ConnectionConfig
	LifecycleBaseURL string
	LifecycleClient  adapters.LifecycleClient
}

type CreateExecdStackOptions struct {
	ConnectionConfig *config.ConnectionConfig
	ExecdBaseURL     string
}

type CreateEgressStackOptions struct {
	ConnectionConfig *config.ConnectionConfig
	EgressBaseURL    string
}

type AdapterFactory interface {
	CreateLifecycleStack(opts CreateLifecycleStackOptions) (*LifecycleStack, error)
	CreateExecdStack(opts CreateExecdStackOptions) (*ExecdStack, error)
	CreateEgressStack(opts CreateEgressStackOptions) (*EgressStack, error)
}
