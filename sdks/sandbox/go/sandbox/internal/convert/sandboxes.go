package convert

import (
	lifecycleapi "github.com/alibaba/opensandbox/sdks/sandbox/go/sandbox/internal/openapi/lifecycle"
	"github.com/alibaba/opensandbox/sdks/sandbox/go/sandbox/models"
)

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

func FromLifecycleCreateSandboxResponse(resp *lifecycleapi.CreateSandboxResponse) *models.CreateSandboxResponse {
	if resp == nil {
		return nil
	}
	return &models.CreateSandboxResponse{
		ID:         resp.Id,
		Status:     fromLifecycleSandboxStatus(resp.Status),
		Metadata:   derefStringMap(resp.Metadata),
		ExpiresAt:  resp.ExpiresAt,
		CreatedAt:  resp.CreatedAt,
		Entrypoint: append([]string(nil), resp.Entrypoint...),
	}
}

func FromLifecycleSandbox(resp *lifecycleapi.Sandbox) *models.SandboxInfo {
	if resp == nil {
		return nil
	}
	return &models.SandboxInfo{
		ID:         resp.Id,
		Image:      fromLifecycleImageSpec(resp.Image),
		Entrypoint: append([]string(nil), resp.Entrypoint...),
		Metadata:   derefStringMap(resp.Metadata),
		Status:     fromLifecycleSandboxStatus(resp.Status),
		CreatedAt:  resp.CreatedAt,
		ExpiresAt:  resp.ExpiresAt,
	}
}

func FromLifecycleListSandboxesResponse(resp *lifecycleapi.ListSandboxesResponse) *models.ListSandboxesResponse {
	if resp == nil {
		return nil
	}
	items := make([]models.SandboxInfo, 0, len(resp.Items))
	for _, item := range resp.Items {
		converted := FromLifecycleSandbox(&item)
		if converted != nil {
			items = append(items, *converted)
		}
	}
	return &models.ListSandboxesResponse{
		Items: items,
		Pagination: &models.PaginationInfo{
			Page:        resp.Pagination.Page,
			PageSize:    resp.Pagination.PageSize,
			TotalItems:  resp.Pagination.TotalItems,
			TotalPages:  resp.Pagination.TotalPages,
			HasNextPage: resp.Pagination.HasNextPage,
		},
	}
}

func FromLifecycleRenewSandboxExpirationResponse(resp *lifecycleapi.RenewSandboxExpirationResponse) *models.RenewSandboxExpirationResponse {
	if resp == nil {
		return nil
	}
	return &models.RenewSandboxExpirationResponse{
		ExpiresAt: &resp.ExpiresAt,
	}
}

func FromLifecycleEndpoint(resp *lifecycleapi.Endpoint) *models.SandboxEndpoint {
	if resp == nil {
		return nil
	}
	return &models.SandboxEndpoint{
		Endpoint: resp.Endpoint,
		Headers:  derefStringMap(resp.Headers),
	}
}

func fromLifecycleSandboxStatus(status lifecycleapi.SandboxStatus) models.SandboxStatus {
	return models.SandboxStatus{
		State:   status.State,
		Reason:  derefString(status.Reason),
		Message: derefString(status.Message),
	}
}

func fromLifecycleImageSpec(image lifecycleapi.ImageSpec) models.ImageSpec {
	out := models.ImageSpec{URI: image.Uri}
	if image.Auth != nil {
		out.Auth = &models.ImageAuth{
			Username: derefString(image.Auth.Username),
			Password: derefString(image.Auth.Password),
		}
	}
	return out
}

func derefStringMap(in *map[string]string) map[string]string {
	if in == nil {
		return nil
	}
	out := make(map[string]string, len(*in))
	for k, v := range *in {
		out[k] = v
	}
	return out
}

func derefString(in *string) string {
	if in == nil {
		return ""
	}
	return *in
}
