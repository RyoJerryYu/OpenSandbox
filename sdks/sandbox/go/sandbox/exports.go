package sandbox

import (
	"github.com/alibaba/opensandbox/sdks/sandbox/go/sandbox/config"
	"github.com/alibaba/opensandbox/sdks/sandbox/go/sandbox/factory"
	"github.com/alibaba/opensandbox/sdks/sandbox/go/sandbox/models"
)

// ConnectionConfig is the public alias for SDK connection settings.
type ConnectionConfig = config.ConnectionConfig

// ImageAuth is the public alias for registry authentication settings.
type ImageAuth = models.ImageAuth
// ImageSpec is the public alias for sandbox image selection.
type ImageSpec = models.ImageSpec
// NetworkRuleAction is the public alias for egress rule actions.
type NetworkRuleAction = models.NetworkRuleAction

const (
	// NetworkRuleActionAllow permits matching outbound traffic.
	NetworkRuleActionAllow = models.NetworkRuleActionAllow
	// NetworkRuleActionDeny blocks matching outbound traffic.
	NetworkRuleActionDeny  = models.NetworkRuleActionDeny
)

// NetworkRule is the public alias for one egress rule.
type NetworkRule = models.NetworkRule
// NetworkPolicy is the public alias for sandbox egress policy configuration.
type NetworkPolicy = models.NetworkPolicy
// Host is the public alias for host-path volume configuration.
type Host = models.Host
// PVC is the public alias for PVC-backed volume configuration.
type PVC = models.PVC
// Volume is the public alias for sandbox volume configuration.
type Volume = models.Volume
// SandboxStatus is the public alias for lifecycle status information.
type SandboxStatus = models.SandboxStatus
// SandboxFilter is the public alias for lifecycle list filters.
type SandboxFilter = models.SandboxFilter
// SandboxInfo is the public alias for lifecycle sandbox information.
type SandboxInfo = models.SandboxInfo
// SandboxEndpoint is the public alias for resolved sandbox service endpoints.
type SandboxEndpoint = models.SandboxEndpoint
// PaginationInfo is the public alias for lifecycle pagination metadata.
type PaginationInfo = models.PaginationInfo
// ListSandboxesResponse is the public alias for paginated sandbox listings.
type ListSandboxesResponse = models.ListSandboxesResponse
// RenewSandboxExpirationResponse is the public alias for renew-expiration results.
type RenewSandboxExpirationResponse = models.RenewSandboxExpirationResponse
// CreateSandboxRequest is the public alias for lifecycle create payloads.
type CreateSandboxRequest = models.CreateSandboxRequest
// CreateSandboxResponse is the public alias for lifecycle create responses.
type CreateSandboxResponse = models.CreateSandboxResponse
// EntryInfo is the public alias for filesystem metadata.
type EntryInfo = models.EntryInfo
// WriteEntry is the public alias for file and directory write requests.
type WriteEntry = models.WriteEntry
// SearchEntry is the public alias for filesystem search requests.
type SearchEntry = models.SearchEntry
// MoveEntry is the public alias for filesystem move requests.
type MoveEntry = models.MoveEntry
// ContentReplaceEntry is the public alias for file content replacement requests.
type ContentReplaceEntry = models.ContentReplaceEntry
// SetPermissionEntry is the public alias for chmod/chown style requests.
type SetPermissionEntry = models.SetPermissionEntry
// ReadFileOptions is the public alias for file download options.
type ReadFileOptions = models.ReadFileOptions
// RunCommandOptions is the public alias for command execution options.
type RunCommandOptions = models.RunCommandOptions
// CommandStatus is the public alias for background command status.
type CommandStatus = models.CommandStatus
// CommandLogs is the public alias for paged background command logs.
type CommandLogs = models.CommandLogs
// OutputMessage is the public alias for one stdout or stderr message.
type OutputMessage = models.OutputMessage
// ExecutionResult is the public alias for one execution result item.
type ExecutionResult = models.ExecutionResult
// ExecutionError is the public alias for execution failure details.
type ExecutionError = models.ExecutionError
// ExecutionLogs is the public alias for grouped stdout and stderr messages.
type ExecutionLogs = models.ExecutionLogs
// ExecutionComplete is the public alias for execution completion timing.
type ExecutionComplete = models.ExecutionComplete
// ServerStreamEvent is the public alias for streamed command events.
type ServerStreamEvent = models.ServerStreamEvent
// CommandExecution is the public alias for aggregated command execution results.
type CommandExecution = models.CommandExecution
// CommandStream is the public alias for streaming command handles.
type CommandStream = models.CommandStream

// AdapterFactory is the public alias for custom adapter stack construction.
type AdapterFactory = factory.AdapterFactory
