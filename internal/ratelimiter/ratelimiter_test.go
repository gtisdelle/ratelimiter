package ratelimiter

import (
	"context"
	"testing"
)

type fakeStore struct {
	allowFunc func() (bool, error)
}

func (s fakeStore) Allow(ctx context.Context, key string) (bool, error) {
	return s.allowFunc()
}

var _ Store = fakeStore{}

func TestAllowUnderLimit(t *testing.T) {
	store := fakeStore{
		allowFunc: func() (bool, error) { return true, nil },
	}
	limiter := NewRateLimiter(store)

	result, err := limiter.Allow(t.Context(), "foo")

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

	result, err := limiter.Allow(t.Context(), "foo")

	if err != nil {
		t.Fatalf("unexpcted error: %v", err)
	}
	if result {
		t.Fatalf("Allow(\"foo\") = %v; want false", result)
	}
}
