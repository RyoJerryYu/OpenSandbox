package models

import "time"

type ImageAuth struct {
	Username string
	Password string
	Token    string
}

type ImageSpec struct {
	URI  string
	Auth *ImageAuth
}

type NetworkRuleAction string

const (
	NetworkRuleActionAllow NetworkRuleAction = "allow"
	NetworkRuleActionDeny  NetworkRuleAction = "deny"
)

type NetworkRule struct {
	Action NetworkRuleAction
	Target string
}

type NetworkPolicy struct {
	DefaultAction NetworkRuleAction
	Egress        []NetworkRule
}

type Host struct {
	Path string
}

type PVC struct {
	ClaimName string
}

type Volume struct {
	Name      string
	Host      *Host
	PVC       *PVC
	MountPath string
	ReadOnly  bool
	SubPath   string
}

type SandboxStatus struct {
	State   string
	Reason  string
	Message string
}

type SandboxInfo struct {
	ID         string
	Image      ImageSpec
	Entrypoint []string
	Metadata   map[string]string
	Status     SandboxStatus
	CreatedAt  time.Time
	ExpiresAt  *time.Time
}

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

type CreateSandboxResponse struct {
	ID         string
	Status     SandboxStatus
	Metadata   map[string]string
	ExpiresAt  *time.Time
	CreatedAt  time.Time
	Entrypoint []string
}

type PaginationInfo struct {
	Page       int
	PageSize   int
	TotalItems int
	TotalPages int
	HasNextPage bool
}

type ListSandboxesResponse struct {
	Items      []SandboxInfo
	Pagination *PaginationInfo
}

type SandboxFilter struct {
	States   []string
	Metadata map[string]string
	Page     int
	PageSize int
}

type RenewSandboxExpirationResponse struct {
	ExpiresAt *time.Time
}

type SandboxEndpoint struct {
	Endpoint string
	Headers  map[string]string
}
