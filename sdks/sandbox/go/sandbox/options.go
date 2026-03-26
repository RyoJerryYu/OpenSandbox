package sandbox

import (
	"context"
	"time"

	"github.com/alibaba/opensandbox/sdks/sandbox/go/sandbox/config"
	"github.com/alibaba/opensandbox/sdks/sandbox/go/sandbox/factory"
	"github.com/alibaba/opensandbox/sdks/sandbox/go/sandbox/models"
)

type SandboxCreateOptions struct {
	ConnectionConfig *config.ConnectionConfig
	AdapterFactory   factory.AdapterFactory
	Timeout          *time.Duration
	Image            models.ImageSpec
	SkipHealthCheck  bool
	ReadyTimeout     time.Duration
	PollInterval     time.Duration
	HealthCheck      func(ctx context.Context, sandbox *Sandbox) (bool, error)
}

type SandboxConnectOptions struct {
	ConnectionConfig *config.ConnectionConfig
	AdapterFactory   factory.AdapterFactory
	SandboxID        string
	SkipHealthCheck  bool
	ReadyTimeout     time.Duration
	PollInterval     time.Duration
	HealthCheck      func(ctx context.Context, sandbox *Sandbox) (bool, error)
}

type ResumeOptions struct {
	ConnectionConfig *config.ConnectionConfig
	AdapterFactory   factory.AdapterFactory
	SkipHealthCheck  bool
	ReadyTimeout     time.Duration
	PollInterval     time.Duration
	HealthCheck      func(ctx context.Context, sandbox *Sandbox) (bool, error)
}
