package factory

import "github.com/alibaba/opensandbox/sdks/sandbox/go/sandbox/services"

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

type CreateLifecycleStackOptions struct{}

type CreateExecdStackOptions struct{}

type CreateEgressStackOptions struct{}

type AdapterFactory interface {
	CreateLifecycleStack(opts CreateLifecycleStackOptions) (*LifecycleStack, error)
	CreateExecdStack(opts CreateExecdStackOptions) (*ExecdStack, error)
	CreateEgressStack(opts CreateEgressStackOptions) (*EgressStack, error)
}
