package ratelimiter

import (
	"context"
	"fmt"
	"strings"

	rlsv3common "github.com/envoyproxy/go-control-plane/envoy/extensions/common/ratelimit/v3"
	rlsv3 "github.com/envoyproxy/go-control-plane/envoy/service/ratelimit/v3"
)

type RateLimiter interface {
	Allow(ctx context.Context, domain string, hits uint64, descriptors []*rlsv3common.RateLimitDescriptor) (*rlsv3.RateLimitResponse, error)
}

type limiter struct {
	store Store
}

func NewRateLimiter(store Store) RateLimiter {
	return &limiter{
		store: store,
	}
}

func (l *limiter) Allow(ctx context.Context, domain string, hits uint64, descriptors []*rlsv3common.RateLimitDescriptor) (*rlsv3.RateLimitResponse, error) {
	statuses := make([]*rlsv3.RateLimitResponse_DescriptorStatus, 0)
	overall := rlsv3.RateLimitResponse_OK
	for _, descriptor := range descriptors {
		allow, err := l.store.Allow(ctx, getKey(domain, descriptor), getHits(hits, descriptor))
		if err != nil {
			return &rlsv3.RateLimitResponse{OverallCode: rlsv3.RateLimitResponse_UNKNOWN}, nil
		}
		code := rlsv3.RateLimitResponse_OK
		if !allow {
			code = rlsv3.RateLimitResponse_OVER_LIMIT
			overall = rlsv3.RateLimitResponse_OVER_LIMIT
		}
		statuses = append(statuses, &rlsv3.RateLimitResponse_DescriptorStatus{
			Code: code,
		})
	}

	return &rlsv3.RateLimitResponse{
		OverallCode: overall,
		Statuses:    statuses,
	}, nil
}

func getKey(domain string, descriptor *rlsv3common.RateLimitDescriptor) string {
	sb := strings.Builder{}
	sb.WriteString(fmt.Sprintf("ratelimit:v1:%s", domain))
	for _, entry := range descriptor.Entries {
		sb.WriteString(":")
		sb.WriteString(fmt.Sprintf("%s:%s", entry.Key, entry.Value))
	}
	return sb.String()
}

func getHits(reqHits uint64, descriptor *rlsv3common.RateLimitDescriptor) uint64 {
	if descriptor.HitsAddend != nil {
		return descriptor.HitsAddend.Value
	}

	return reqHits
}
