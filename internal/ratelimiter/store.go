package ratelimiter

import (
	"context"
	"time"
)

type Store interface {
	Increment(ctx context.Context, key string, ttl time.Duration) (int, error)
}
