package sandbox

import (
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
}

type SandboxConnectOptions struct {
	ConnectionConfig *config.ConnectionConfig
	AdapterFactory   factory.AdapterFactory
	SandboxID        string
}

type ResumeOptions struct {
	ConnectionConfig *config.ConnectionConfig
	AdapterFactory   factory.AdapterFactory
}
