package convert

import "github.com/alibaba/opensandbox/sdks/sandbox/go/sandbox/models"

func CloneCreateSandboxResponse(resp *models.CreateSandboxResponse) *models.CreateSandboxResponse {
	if resp == nil {
		return nil
	}
	clone := *resp
	return &clone
}

func CloneSandboxInfo(resp *models.SandboxInfo) *models.SandboxInfo {
	if resp == nil {
		return nil
	}
	clone := *resp
	return &clone
}

func CloneListSandboxesResponse(resp *models.ListSandboxesResponse) *models.ListSandboxesResponse {
	if resp == nil {
		return nil
	}
	clone := *resp
	return &clone
}

func CloneRenewSandboxExpirationResponse(resp *models.RenewSandboxExpirationResponse) *models.RenewSandboxExpirationResponse {
	if resp == nil {
		return nil
	}
	clone := *resp
	return &clone
}

func CloneSandboxEndpoint(resp *models.SandboxEndpoint) *models.SandboxEndpoint {
	if resp == nil {
		return nil
	}
	clone := *resp
	return &clone
}
