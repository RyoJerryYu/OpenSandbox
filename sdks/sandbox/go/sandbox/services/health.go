package services

import "context"

// ExecdHealth defines health checks against the execd sidecar.
type ExecdHealth interface {
	Ping(ctx context.Context) (bool, error)
}
