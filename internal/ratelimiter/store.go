package ratelimiter

import (
	"context"
)

type Store interface {
	Allow(ctx context.Context, key limitKey) (bool, error)
}
