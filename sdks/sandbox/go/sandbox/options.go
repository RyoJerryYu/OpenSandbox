package sandbox

import (
	"context"
	"time"

	"github.com/alibaba/opensandbox/sdks/sandbox/go/sandbox/config"
	"github.com/alibaba/opensandbox/sdks/sandbox/go/sandbox/factory"
	"github.com/alibaba/opensandbox/sdks/sandbox/go/sandbox/models"
)

// SandboxCreateOptions controls sandbox creation and default readiness behavior.
type SandboxCreateOptions struct {
	// ConnectionConfig configures lifecycle and sandbox endpoint access.
	ConnectionConfig *config.ConnectionConfig
	// AdapterFactory overrides the default generated transport wiring.
	AdapterFactory factory.AdapterFactory
	// Timeout is converted to lifecycle timeout seconds. Nil leaves the server default in effect.
	Timeout *time.Duration
	// Image selects the runtime image to boot.
	Image models.ImageSpec
	// SkipHealthCheck disables default readiness waiting after create.
	SkipHealthCheck bool
	// ReadyTimeout bounds readiness polling when SkipHealthCheck is false.
	ReadyTimeout time.Duration
	// PollInterval controls readiness poll frequency.
	PollInterval time.Duration
	// HealthCheck replaces the default execd ping during readiness polling.
	HealthCheck func(ctx context.Context, sandbox Sandbox) (bool, error)
}

// SandboxConnectOptions controls attachment to an existing sandbox.
type SandboxConnectOptions struct {
	ConnectionConfig *config.ConnectionConfig
	AdapterFactory   factory.AdapterFactory
	SandboxID        string
	SkipHealthCheck  bool
	ReadyTimeout     time.Duration
	PollInterval     time.Duration
	HealthCheck      func(ctx context.Context, sandbox Sandbox) (bool, error)
}

// ResumeOptions controls reconnection behavior after resuming a paused sandbox.
type ResumeOptions struct {
	ConnectionConfig *config.ConnectionConfig
	AdapterFactory   factory.AdapterFactory
	SkipHealthCheck  bool
	ReadyTimeout     time.Duration
	PollInterval     time.Duration
	HealthCheck      func(ctx context.Context, sandbox Sandbox) (bool, error)
}
