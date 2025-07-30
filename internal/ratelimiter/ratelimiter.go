package ratelimiter

import (
	"fmt"
	"time"
)

type RateLimiter interface {
	Allow(key string) (bool, error)
}

type limiter struct {
	store      Store
	clock      Clock
	limit      int
	windowSize time.Duration
}

func NewRateLimiter(store Store, clock Clock, limit int, windowSize time.Duration) RateLimiter {
	return &limiter{
		store:      store,
		clock:      clock,
		limit:      limit,
		windowSize: windowSize,
	}
}

func (l *limiter) Allow(key string) (bool, error) {
	windowStart := l.clock.Now().UnixMilli() - (l.clock.Now().UnixMilli() % l.windowSize.Milliseconds())
	cacheKey := fmt.Sprintf("ratelimiter:v1:%v:%v", key, windowStart)
	count, err := l.store.Increment(cacheKey, int(l.windowSize.Milliseconds()))
	if err != nil {
		return false, err
	}

	return count <= l.limit, nil
}
