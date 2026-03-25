package services

import (
	"context"

	"github.com/alibaba/opensandbox/sdks/sandbox/go/sandbox/models"
)

type Egress interface {
	GetPolicy(ctx context.Context) (*models.NetworkPolicy, error)
	PatchRules(ctx context.Context, rules []models.NetworkRule) error
}
