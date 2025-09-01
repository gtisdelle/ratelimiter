package ratelimiter

import (
	"context"
)

type Store interface {
	Allow(ctx context.Context, key string, hits uint64) (bool, error)
}
