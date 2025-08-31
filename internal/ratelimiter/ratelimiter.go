package ratelimiter

import (
	"context"

	rlsv3common "github.com/envoyproxy/go-control-plane/envoy/extensions/common/ratelimit/v3"
)

type RateLimiter interface {
	Allow(ctx context.Context, domain string, hits uint64, descriptors []*rlsv3common.RateLimitDescriptor) (bool, error)
}

type limiter struct {
	store Store
}

func NewRateLimiter(store Store) RateLimiter {
	return &limiter{
		store: store,
	}
}

func (l *limiter) Allow(ctx context.Context, domain string, hits uint64, descriptors []*rlsv3common.RateLimitDescriptor) (bool, error) {
	for _, descriptor := range descriptors {
		key := NewLimitKey(domain, hits, descriptor)
		allow, err := l.store.Allow(ctx, key)
		if err != nil {
			return false, err
		}
		if !allow {
			return false, nil
		}
	}
	return true, nil
}
