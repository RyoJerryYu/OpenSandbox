package adapters

import (
	"context"
	"fmt"

	"github.com/alibaba/opensandbox/sdks/sandbox/go/sandbox/internal/convert"
	egressapi "github.com/alibaba/opensandbox/sdks/sandbox/go/sandbox/internal/openapi/egress"
	"github.com/alibaba/opensandbox/sdks/sandbox/go/sandbox/models"
)

type EgressClient interface {
	GetPolicyWithResponse(ctx context.Context, reqEditors ...egressapi.RequestEditorFn) (*egressapi.GetPolicyResponse, error)
	PatchPolicyWithResponse(ctx context.Context, body egressapi.PatchPolicyJSONRequestBody, reqEditors ...egressapi.RequestEditorFn) (*egressapi.PatchPolicyResponse, error)
}

type EgressAdapter struct {
	client EgressClient
}

func NewEgressAdapter(client EgressClient) *EgressAdapter {
	return &EgressAdapter{client: client}
}

func (a *EgressAdapter) GetPolicy(ctx context.Context) (*models.NetworkPolicy, error) {
	resp, err := a.client.GetPolicyWithResponse(ctx)
	if err != nil {
		return nil, err
	}
	if resp.JSON200 == nil {
		return nil, NormalizeHTTPError(fmt.Errorf("get egress policy failed: unexpected status %d", resp.StatusCode()), resp.HTTPResponse)
	}
	return convert.FromEgressPolicy(resp.JSON200), nil
}

func (a *EgressAdapter) PatchRules(ctx context.Context, rules []models.NetworkRule) error {
	resp, err := a.client.PatchPolicyWithResponse(ctx, convert.ToEgressRules(rules))
	if err != nil {
		return err
	}
	if resp.StatusCode() >= 200 && resp.StatusCode() < 300 {
		return nil
	}
	return NormalizeHTTPError(fmt.Errorf("patch egress rules failed: unexpected status %d", resp.StatusCode()), resp.HTTPResponse)
}
