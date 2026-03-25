package convert

import (
	egressapi "github.com/alibaba/opensandbox/sdks/sandbox/go/sandbox/internal/openapi/egress"
	"github.com/alibaba/opensandbox/sdks/sandbox/go/sandbox/models"
)

func FromEgressPolicy(resp *egressapi.PolicyStatusResponse) *models.NetworkPolicy {
	if resp == nil || resp.Policy == nil {
		return &models.NetworkPolicy{}
	}

	policy := &models.NetworkPolicy{}
	if resp.Policy.DefaultAction != nil {
		policy.DefaultAction = models.NetworkRuleAction(*resp.Policy.DefaultAction)
	}
	if resp.Policy.Egress != nil {
		policy.Egress = make([]models.NetworkRule, 0, len(*resp.Policy.Egress))
		for _, rule := range *resp.Policy.Egress {
			policy.Egress = append(policy.Egress, models.NetworkRule{
				Action: models.NetworkRuleAction(rule.Action),
				Target: rule.Target,
			})
		}
	}
	return policy
}

func ToEgressRules(rules []models.NetworkRule) egressapi.PatchPolicyJSONRequestBody {
	out := make(egressapi.PatchPolicyJSONRequestBody, 0, len(rules))
	for _, rule := range rules {
		out = append(out, egressapi.NetworkRule{
			Action: egressapi.NetworkRuleAction(rule.Action),
			Target: rule.Target,
		})
	}
	return out
}
