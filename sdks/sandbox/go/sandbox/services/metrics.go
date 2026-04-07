package services

import (
	"context"

	"github.com/alibaba/opensandbox/sdks/sandbox/go/sandbox/models"
)

// ExecdMetrics defines resource metric queries against execd.
type ExecdMetrics interface {
	GetMetrics(ctx context.Context) (*models.SandboxMetrics, error)
}
