package limit

import (
	"context"
	"fmt"

	rlsv3common "github.com/envoyproxy/go-control-plane/envoy/extensions/common/ratelimit/v3"
	rlsv3 "github.com/envoyproxy/go-control-plane/envoy/service/ratelimit/v3"
)

var (
	DEFAULT = "DEFAULT"
)

type store interface {
	Allow(ctx context.Context, key string, hits uint64) (bool, int, error)
}

type Limiter struct {
	store        store
	defaultLimit int
}

func NewLimiter(store store, defaultLimit int) *Limiter {
	return &Limiter{
		store:        store,
		defaultLimit: defaultLimit,
	}
}

func (l *Limiter) Allow(ctx context.Context, domain string, hits uint64, descriptors []*rlsv3common.RateLimitDescriptor) (*rlsv3.RateLimitResponse, error) {
	statuses := make([]*rlsv3.RateLimitResponse_DescriptorStatus, 0)
	overall := rlsv3.RateLimitResponse_OK
	for _, descriptor := range descriptors {
		key := BuildKey(domain, descriptor)
		allow, remaining, err := l.store.Allow(ctx, key, getHits(hits, descriptor))
		if err != nil {
			return nil, fmt.Errorf("limit %s: %w", key, err)
		}
		code := rlsv3.RateLimitResponse_OK
		if !allow {
			code = rlsv3.RateLimitResponse_OVER_LIMIT
			overall = rlsv3.RateLimitResponse_OVER_LIMIT
		}
		statuses = append(statuses, &rlsv3.RateLimitResponse_DescriptorStatus{
			Code: code,
			CurrentLimit: &rlsv3.RateLimitResponse_RateLimit{
				Name:            DEFAULT,
				RequestsPerUnit: uint32(l.defaultLimit),
				Unit:            rlsv3.RateLimitResponse_RateLimit_SECOND,
			},
			LimitRemaining: uint32(remaining),
		})
	}

	return &rlsv3.RateLimitResponse{
		OverallCode: overall,
		Statuses:    statuses,
	}, nil
}

func getHits(reqHits uint64, descriptor *rlsv3common.RateLimitDescriptor) uint64 {
	if descriptor.HitsAddend != nil {
		return descriptor.HitsAddend.Value
	}

	return reqHits
}
