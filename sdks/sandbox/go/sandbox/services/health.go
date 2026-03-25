package services

import "context"

type ExecdHealth interface {
	Ping(ctx context.Context) (bool, error)
}
