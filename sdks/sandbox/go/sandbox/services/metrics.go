package services

import (
	"context"

	"github.com/alibaba/opensandbox/sdks/sandbox/go/sandbox/models"
)

type ExecdMetrics interface {
	GetMetrics(ctx context.Context) (*models.SandboxMetrics, error)
}
