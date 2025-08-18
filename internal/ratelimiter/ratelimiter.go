package ratelimiter

import (
	"context"
	"fmt"
)

type RateLimiter interface {
	Allow(ctx context.Context, key string) (bool, error)
}

type limiter struct {
	store Store
}

func NewRateLimiter(store Store) RateLimiter {
	return &limiter{
		store: store,
	}
}

func (l *limiter) Allow(ctx context.Context, key string) (bool, error) {
	cacheKey := fmt.Sprintf("ratelimiter:v1:%v", key)
	allow, err := l.store.Allow(ctx, cacheKey)
	if err != nil {
		return false, err
	}
	return allow, nil
}
