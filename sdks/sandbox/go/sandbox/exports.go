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
type EntryInfo = models.EntryInfo
type WriteEntry = models.WriteEntry
type SearchEntry = models.SearchEntry
type MoveEntry = models.MoveEntry
type ContentReplaceEntry = models.ContentReplaceEntry
type SetPermissionEntry = models.SetPermissionEntry
type ReadFileOptions = models.ReadFileOptions
type RunCommandOptions = models.RunCommandOptions
type CommandStatus = models.CommandStatus
type CommandLogs = models.CommandLogs

type AdapterFactory = factory.AdapterFactory
