package sandbox

import (
	"github.com/alibaba/opensandbox/sdks/sandbox/go/sandbox/config"
	"github.com/alibaba/opensandbox/sdks/sandbox/go/sandbox/factory"
	"github.com/alibaba/opensandbox/sdks/sandbox/go/sandbox/models"
)

type ConnectionConfig = config.ConnectionConfig

type SandboxFilter = models.SandboxFilter
type SandboxInfo = models.SandboxInfo
type SandboxEndpoint = models.SandboxEndpoint
type ListSandboxesResponse = models.ListSandboxesResponse
type RenewSandboxExpirationResponse = models.RenewSandboxExpirationResponse
type CreateSandboxRequest = models.CreateSandboxRequest
type CreateSandboxResponse = models.CreateSandboxResponse

type AdapterFactory = factory.AdapterFactory
