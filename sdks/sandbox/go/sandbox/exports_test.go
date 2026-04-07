package sandbox

import "testing"

func TestPublicExportsExposeCommonSandboxModels(t *testing.T) {
	req := CreateSandboxRequest{
		Image: ImageSpec{
			URI: "ubuntu",
			Auth: &ImageAuth{
				Username: "demo",
				Password: "secret",
			},
		},
		NetworkPolicy: &NetworkPolicy{
			DefaultAction: NetworkRuleActionDeny,
			Egress: []NetworkRule{
				{Action: NetworkRuleActionAllow, Target: "pypi.org"},
			},
		},
		Volumes: []Volume{
			{
				Name:      "workspace",
				MountPath: "/workspace",
				Host:      &Host{Path: "/tmp/workspace"},
			},
			{
				Name:      "cache",
				MountPath: "/cache",
				PVC:       &PVC{ClaimName: "sandbox-cache"},
			},
		},
	}

	if req.Image.URI != "ubuntu" {
		t.Fatalf("unexpected image uri: %s", req.Image.URI)
	}
	if req.NetworkPolicy == nil || req.NetworkPolicy.DefaultAction != NetworkRuleActionDeny {
		t.Fatalf("unexpected network policy: %+v", req.NetworkPolicy)
	}
	if len(req.Volumes) != 2 {
		t.Fatalf("unexpected volume count: %d", len(req.Volumes))
	}
}
