package sandbox

import (
	"time"

	"github.com/alibaba/opensandbox/sdks/sandbox/go/sandbox/config"
)

type SandboxCreateOptions struct {
	ConnectionConfig *config.ConnectionConfig
	Timeout          *time.Duration
}
