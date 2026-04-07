package models

import "time"

type ImageAuth struct {
	// Username is the registry username for private image pulls.
	Username string
	// Password is the registry password or access token for private image pulls.
	Password string
	// Token carries bearer-style auth when the registry uses token authentication.
	Token    string
}

// ImageSpec identifies the sandbox image to run.
type ImageSpec struct {
	// URI is the OCI image reference, for example "ubuntu:24.04".
	URI  string
	// Auth contains optional registry credentials for private images.
	Auth *ImageAuth
}

// NetworkRuleAction is the action applied when a network rule matches.
type NetworkRuleAction string

const (
	NetworkRuleActionAllow NetworkRuleAction = "allow"
	NetworkRuleActionDeny  NetworkRuleAction = "deny"
)

// NetworkRule matches outbound traffic against a target and applies an action.
type NetworkRule struct {
	Action NetworkRuleAction
	Target string
}

// NetworkPolicy defines outbound network behavior for a sandbox.
type NetworkPolicy struct {
	DefaultAction NetworkRuleAction
	Egress        []NetworkRule
}

// Host configures a host-path backed volume.
type Host struct {
	Path string
}

// PVC configures a Kubernetes PersistentVolumeClaim backed volume.
type PVC struct {
	ClaimName string
}

// Volume describes a filesystem mount attached to the sandbox.
type Volume struct {
	Name      string
	Host      *Host
	PVC       *PVC
	MountPath string
	ReadOnly  bool
	SubPath   string
}

// SandboxStatus is the lifecycle status reported by the server.
type SandboxStatus struct {
	State   string
	Reason  string
	Message string
}

// SandboxInfo is the lifecycle view of a sandbox instance.
type SandboxInfo struct {
	ID         string
	Image      ImageSpec
	Entrypoint []string
	Metadata   map[string]string
	Status     SandboxStatus
	CreatedAt  time.Time
	ExpiresAt  *time.Time
}

// CreateSandboxRequest contains the lifecycle payload for creating a sandbox.
type CreateSandboxRequest struct {
	Image         ImageSpec
	Entrypoint    []string
	Timeout       *int
	ResourceLimits map[string]string
	Env           map[string]string
	Metadata      map[string]string
	NetworkPolicy *NetworkPolicy
	Volumes       []Volume
	Extensions    map[string]string
}

// CreateSandboxResponse contains the initial lifecycle response after sandbox creation.
type CreateSandboxResponse struct {
	ID         string
	Status     SandboxStatus
	Metadata   map[string]string
	ExpiresAt  *time.Time
	CreatedAt  time.Time
	Entrypoint []string
}

// PaginationInfo describes paginated lifecycle list results.
type PaginationInfo struct {
	Page       int
	PageSize   int
	TotalItems int
	TotalPages int
	HasNextPage bool
}

// ListSandboxesResponse contains lifecycle list results and optional pagination metadata.
type ListSandboxesResponse struct {
	Items      []SandboxInfo
	Pagination *PaginationInfo
}

// SandboxFilter narrows lifecycle list queries by state, metadata, and pagination.
type SandboxFilter struct {
	States   []string
	Metadata map[string]string
	Page     int
	PageSize int
}

// RenewSandboxExpirationResponse contains the updated expiration time after a renew request.
type RenewSandboxExpirationResponse struct {
	ExpiresAt *time.Time
}

// SandboxEndpoint describes how to reach a service exposed from the sandbox.
type SandboxEndpoint struct {
	// Endpoint is a host[:port][/path] value without scheme.
	Endpoint string
	// Headers contains extra request headers required by the endpoint, if any.
	Headers  map[string]string
}
