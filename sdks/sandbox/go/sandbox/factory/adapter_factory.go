package factory

import (
	"github.com/alibaba/opensandbox/sdks/sandbox/go/sandbox/adapters"
	"github.com/alibaba/opensandbox/sdks/sandbox/go/sandbox/config"
	"github.com/alibaba/opensandbox/sdks/sandbox/go/sandbox/services"
)

// LifecycleStack groups lifecycle-oriented service adapters.
type LifecycleStack struct {
	Sandboxes services.Sandboxes
}

// ExecdStack groups execd-oriented service adapters.
type ExecdStack struct {
	Commands services.ExecdCommands
	Files    services.SandboxFiles
	Health   services.ExecdHealth
	Metrics  services.ExecdMetrics
}

// EgressStack groups egress sidecar adapters.
type EgressStack struct {
	Egress services.Egress
}

// CreateLifecycleStackOptions configures lifecycle stack construction.
type CreateLifecycleStackOptions struct {
	ConnectionConfig *config.ConnectionConfig
	LifecycleBaseURL string
	LifecycleClient  adapters.LifecycleClient
}

// CreateExecdStackOptions configures execd stack construction.
type CreateExecdStackOptions struct {
	ConnectionConfig *config.ConnectionConfig
	ExecdBaseURL     string
}

// CreateEgressStackOptions configures egress stack construction.
type CreateEgressStackOptions struct {
	ConnectionConfig *config.ConnectionConfig
	EgressBaseURL    string
}

// AdapterFactory builds lifecycle and per-sandbox adapter stacks.
type AdapterFactory interface {
	CreateLifecycleStack(opts CreateLifecycleStackOptions) (*LifecycleStack, error)
	CreateExecdStack(opts CreateExecdStackOptions) (*ExecdStack, error)
	CreateEgressStack(opts CreateEgressStackOptions) (*EgressStack, error)
}
