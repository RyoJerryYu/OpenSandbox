package adapters

import (
	"context"
	"net/http"
	"testing"

	egressapi "github.com/alibaba/opensandbox/sdks/sandbox/go/sandbox/internal/openapi/egress"
	"github.com/alibaba/opensandbox/sdks/sandbox/go/sandbox/models"
)

func TestEgressAdapterGetsPolicy(t *testing.T) {
	adapter := NewEgressAdapter(&fakeEgressClient{
		getPolicyResponse: &egressapi.GetPolicyResponse{
			HTTPResponse: &http.Response{StatusCode: http.StatusOK},
			JSON200: &egressapi.PolicyStatusResponse{
				Policy: &egressapi.NetworkPolicy{
					DefaultAction: egressPolicyActionPtr(egressapi.NetworkPolicyDefaultActionDeny),
					Egress: &[]egressapi.NetworkRule{
						{Action: egressapi.NetworkRuleActionAllow, Target: "pypi.org"},
					},
				},
			},
		},
	})

	policy, err := adapter.GetPolicy(context.Background())
	if err != nil {
		t.Fatalf("get policy: %v", err)
	}
	if policy.DefaultAction != models.NetworkRuleActionDeny {
		t.Fatalf("unexpected default action: %s", policy.DefaultAction)
	}
	if len(policy.Egress) != 1 || policy.Egress[0].Target != "pypi.org" {
		t.Fatalf("unexpected egress rules: %+v", policy.Egress)
	}
}

func TestEgressAdapterPatchesRulesWithoutReplacingDefaultAction(t *testing.T) {
	client := &fakeEgressClient{
		patchPolicyResponse: &egressapi.PatchPolicyResponse{
			HTTPResponse: &http.Response{StatusCode: http.StatusOK},
			JSON200:      &egressapi.PolicyStatusResponse{},
		},
	}
	adapter := NewEgressAdapter(client)

	err := adapter.PatchRules(context.Background(), []models.NetworkRule{
		{Action: models.NetworkRuleActionAllow, Target: "github.com"},
		{Action: models.NetworkRuleActionDeny, Target: "example.com"},
	})
	if err != nil {
		t.Fatalf("patch rules: %v", err)
	}
	if len(client.patchBody) != 2 {
		t.Fatalf("unexpected patch body length: %d", len(client.patchBody))
	}
	if client.patchBody[0].Target != "github.com" || client.patchBody[0].Action != egressapi.NetworkRuleActionAllow {
		t.Fatalf("unexpected first patch rule: %+v", client.patchBody[0])
	}
}

type fakeEgressClient struct {
	getPolicyResponse   *egressapi.GetPolicyResponse
	patchPolicyResponse *egressapi.PatchPolicyResponse
	patchBody           egressapi.PatchPolicyJSONRequestBody
}

func (f *fakeEgressClient) GetPolicyWithResponse(ctx context.Context, reqEditors ...egressapi.RequestEditorFn) (*egressapi.GetPolicyResponse, error) {
	return f.getPolicyResponse, nil
}

func (f *fakeEgressClient) PatchPolicyWithResponse(ctx context.Context, body egressapi.PatchPolicyJSONRequestBody, reqEditors ...egressapi.RequestEditorFn) (*egressapi.PatchPolicyResponse, error) {
	f.patchBody = body
	return f.patchPolicyResponse, nil
}

func egressPolicyActionPtr(v egressapi.NetworkPolicyDefaultAction) *egressapi.NetworkPolicyDefaultAction {
	return &v
}
