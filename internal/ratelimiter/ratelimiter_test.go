package ratelimiter

import (
	"context"
	"testing"

	ratelimitv3 "github.com/envoyproxy/go-control-plane/envoy/extensions/common/ratelimit/v3"
)

type fakeStore struct {
	allowFunc func() (bool, error)
}

func (s fakeStore) Allow(ctx context.Context, key limitKey) (bool, error) {
	return s.allowFunc()
}

var _ Store = fakeStore{}

func TestAllowUnderLimit(t *testing.T) {
	store := fakeStore{
		allowFunc: func() (bool, error) { return true, nil },
	}
	limiter := NewRateLimiter(store)

	result, err := limiter.Allow(t.Context(), "foo", 1, make([]*ratelimitv3.RateLimitDescriptor, 0))

	if err != nil {
		t.Fatalf("unexpcted error: %v", err)
	}
	if !result {
		t.Fatalf("Allow(\"foo\") = %v; want true", result)
	}
}

func TestAllowOverLimit(t *testing.T) {
	store := fakeStore{
		allowFunc: func() (bool, error) { return false, nil },
	}
	limiter := NewRateLimiter(store)
	descriptors := []*ratelimitv3.RateLimitDescriptor{
		{Entries: []*ratelimitv3.RateLimitDescriptor_Entry{{Key: "type", Value: "legacy"}}}}
	result, err := limiter.Allow(t.Context(), "foo", 1, descriptors)

	if err != nil {
		t.Fatalf("unexpcted error: %v", err)
	}
	if result {
		t.Fatalf("Allow(\"foo\") = %v; want false", result)
	}
}
