package adapters

import (
	"context"
	"fmt"
	"net/url"
	"sort"
	"strings"
	"time"

	"github.com/alibaba/opensandbox/sdks/sandbox/go/sandbox/internal/convert"
	lifecycleapi "github.com/alibaba/opensandbox/sdks/sandbox/go/sandbox/internal/openapi/lifecycle"
	"github.com/alibaba/opensandbox/sdks/sandbox/go/sandbox/models"
)

type GeneratedLifecycleClient struct {
	client lifecycleapi.ClientWithResponsesInterface
}

func NewGeneratedLifecycleClient(client lifecycleapi.ClientWithResponsesInterface) *GeneratedLifecycleClient {
	return &GeneratedLifecycleClient{client: client}
}

func (c *GeneratedLifecycleClient) CreateSandbox(ctx context.Context, req models.CreateSandboxRequest) (*models.CreateSandboxResponse, error) {
	resp, err := c.client.PostSandboxesWithResponse(ctx, toLifecycleCreateSandboxRequest(req))
	if err != nil {
		return nil, err
	}
	if resp.JSON202 == nil {
		return nil, NormalizeHTTPError(fmt.Errorf("create sandbox failed: unexpected status %d", resp.StatusCode()), resp.HTTPResponse)
	}
	return convert.FromLifecycleCreateSandboxResponse(resp.JSON202), nil
}

func (c *GeneratedLifecycleClient) GetSandbox(ctx context.Context, sandboxID string) (*models.SandboxInfo, error) {
	resp, err := c.client.GetSandboxesSandboxIdWithResponse(ctx, sandboxID)
	if err != nil {
		return nil, err
	}
	if resp.JSON200 == nil {
		return nil, NormalizeHTTPError(fmt.Errorf("get sandbox failed: unexpected status %d", resp.StatusCode()), resp.HTTPResponse)
	}
	return convert.FromLifecycleSandbox(resp.JSON200), nil
}

func (c *GeneratedLifecycleClient) ListSandboxes(ctx context.Context, filter models.SandboxFilter) (*models.ListSandboxesResponse, error) {
	params := &lifecycleapi.GetSandboxesParams{}
	if len(filter.States) > 0 {
		states := append([]string(nil), filter.States...)
		params.State = &states
	}
	if len(filter.Metadata) > 0 {
		encoded := encodeMetadataFilter(filter.Metadata)
		params.Metadata = &encoded
	}
	if filter.Page > 0 {
		params.Page = &filter.Page
	}
	if filter.PageSize > 0 {
		params.PageSize = &filter.PageSize
	}

	resp, err := c.client.GetSandboxesWithResponse(ctx, params)
	if err != nil {
		return nil, err
	}
	if resp.JSON200 == nil {
		return nil, NormalizeHTTPError(fmt.Errorf("list sandboxes failed: unexpected status %d", resp.StatusCode()), resp.HTTPResponse)
	}
	return convert.FromLifecycleListSandboxesResponse(resp.JSON200), nil
}

func (c *GeneratedLifecycleClient) DeleteSandbox(ctx context.Context, sandboxID string) error {
	resp, err := c.client.DeleteSandboxesSandboxIdWithResponse(ctx, sandboxID)
	if err != nil {
		return err
	}
	if resp.StatusCode() >= 200 && resp.StatusCode() < 300 {
		return nil
	}
	return NormalizeHTTPError(fmt.Errorf("delete sandbox failed: unexpected status %d", resp.StatusCode()), resp.HTTPResponse)
}

func (c *GeneratedLifecycleClient) PauseSandbox(ctx context.Context, sandboxID string) error {
	resp, err := c.client.PostSandboxesSandboxIdPauseWithResponse(ctx, sandboxID)
	if err != nil {
		return err
	}
	if resp.StatusCode() >= 200 && resp.StatusCode() < 300 {
		return nil
	}
	return NormalizeHTTPError(fmt.Errorf("pause sandbox failed: unexpected status %d", resp.StatusCode()), resp.HTTPResponse)
}

func (c *GeneratedLifecycleClient) ResumeSandbox(ctx context.Context, sandboxID string) error {
	resp, err := c.client.PostSandboxesSandboxIdResumeWithResponse(ctx, sandboxID)
	if err != nil {
		return err
	}
	if resp.StatusCode() >= 200 && resp.StatusCode() < 300 {
		return nil
	}
	return NormalizeHTTPError(fmt.Errorf("resume sandbox failed: unexpected status %d", resp.StatusCode()), resp.HTTPResponse)
}

func (c *GeneratedLifecycleClient) RenewSandboxExpiration(ctx context.Context, sandboxID string, expiresAt time.Time) (*models.RenewSandboxExpirationResponse, error) {
	resp, err := c.client.PostSandboxesSandboxIdRenewExpirationWithResponse(ctx, sandboxID, lifecycleapi.RenewSandboxExpirationRequest{
		ExpiresAt: expiresAt,
	})
	if err != nil {
		return nil, err
	}
	if resp.JSON200 == nil {
		return nil, NormalizeHTTPError(fmt.Errorf("renew sandbox failed: unexpected status %d", resp.StatusCode()), resp.HTTPResponse)
	}
	return convert.FromLifecycleRenewSandboxExpirationResponse(resp.JSON200), nil
}

func (c *GeneratedLifecycleClient) GetSandboxEndpoint(ctx context.Context, sandboxID string, port int, useServerProxy bool) (*models.SandboxEndpoint, error) {
	params := &lifecycleapi.GetSandboxesSandboxIdEndpointsPortParams{
		UseServerProxy: &useServerProxy,
	}
	resp, err := c.client.GetSandboxesSandboxIdEndpointsPortWithResponse(ctx, sandboxID, port, params)
	if err != nil {
		return nil, err
	}
	if resp.JSON200 == nil {
		return nil, NormalizeHTTPError(fmt.Errorf("get sandbox endpoint failed: unexpected status %d", resp.StatusCode()), resp.HTTPResponse)
	}
	return convert.FromLifecycleEndpoint(resp.JSON200), nil
}

func toLifecycleCreateSandboxRequest(req models.CreateSandboxRequest) lifecycleapi.CreateSandboxRequest {
	out := lifecycleapi.CreateSandboxRequest{
		Image: lifecycleapi.ImageSpec{
			Uri: req.Image.URI,
		},
		Entrypoint:     append([]string(nil), req.Entrypoint...),
		ResourceLimits: lifecycleapi.ResourceLimits(cloneStringMap(req.ResourceLimits)),
	}
	if req.Image.Auth != nil {
		out.Image.Auth = &struct {
			Password *string `json:"password,omitempty"`
			Username *string `json:"username,omitempty"`
		}{
			Username: stringPtr(req.Image.Auth.Username),
			Password: stringPtr(req.Image.Auth.Password),
		}
	}
	if req.Timeout != nil {
		out.Timeout = req.Timeout
	}
	if req.Env != nil {
		env := cloneStringMap(req.Env)
		out.Env = &env
	}
	if req.Metadata != nil {
		metadata := cloneStringMap(req.Metadata)
		out.Metadata = &metadata
	}
	if req.Extensions != nil {
		extensions := cloneStringMap(req.Extensions)
		out.Extensions = &extensions
	}
	if req.NetworkPolicy != nil {
		out.NetworkPolicy = toLifecycleNetworkPolicy(req.NetworkPolicy)
	}
	if len(req.Volumes) > 0 {
		volumes := make([]lifecycleapi.Volume, 0, len(req.Volumes))
		for _, volume := range req.Volumes {
			volumes = append(volumes, toLifecycleVolume(volume))
		}
		out.Volumes = &volumes
	}
	return out
}

func toLifecycleNetworkPolicy(policy *models.NetworkPolicy) *lifecycleapi.NetworkPolicy {
	if policy == nil {
		return nil
	}
	out := &lifecycleapi.NetworkPolicy{}
	if policy.DefaultAction != "" {
		action := lifecycleapi.NetworkPolicyDefaultAction(policy.DefaultAction)
		out.DefaultAction = &action
	}
	if len(policy.Egress) > 0 {
		egress := make([]lifecycleapi.NetworkRule, 0, len(policy.Egress))
		for _, rule := range policy.Egress {
			egress = append(egress, lifecycleapi.NetworkRule{
				Action: lifecycleapi.NetworkRuleAction(rule.Action),
				Target: rule.Target,
			})
		}
		out.Egress = &egress
	}
	return out
}

func toLifecycleVolume(volume models.Volume) lifecycleapi.Volume {
	out := lifecycleapi.Volume{
		Name:      volume.Name,
		MountPath: volume.MountPath,
	}
	if volume.Host != nil {
		out.Host = &lifecycleapi.Host{Path: volume.Host.Path}
	}
	if volume.PVC != nil {
		out.Pvc = &lifecycleapi.PVC{ClaimName: volume.PVC.ClaimName}
	}
	out.ReadOnly = boolPtr(volume.ReadOnly)
	if volume.SubPath != "" {
		out.SubPath = stringPtr(volume.SubPath)
	}
	return out
}

func encodeMetadataFilter(metadata map[string]string) string {
	if len(metadata) == 0 {
		return ""
	}
	keys := make([]string, 0, len(metadata))
	for key := range metadata {
		keys = append(keys, key)
	}
	sort.Strings(keys)
	parts := make([]string, 0, len(keys))
	for _, key := range keys {
		parts = append(parts, url.QueryEscape(key)+"="+url.QueryEscape(metadata[key]))
	}
	return strings.Join(parts, "&")
}

func cloneStringMap(in map[string]string) map[string]string {
	if in == nil {
		return nil
	}
	out := make(map[string]string, len(in))
	for k, v := range in {
		out[k] = v
	}
	return out
}

func stringPtr(v string) *string { return &v }

func boolPtr(v bool) *bool { return &v }
