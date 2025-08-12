package ratelimiter

import (
	"context"
	"testing"
	"time"
)

type fakeStore struct {
	incrFunc func() (int, error)
}

func (s fakeStore) Increment(ctx context.Context, key string, ttl time.Duration) (int, error) {
	return s.incrFunc()
}

func TestAllowUnderLimit(t *testing.T) {
	clock := &MockClock{currentTime: time.Date(2000, 1, 1, 0, 0, 0, 0, time.UTC)}
	store := fakeStore{
		incrFunc: func() (int, error) { return 1, nil },
	}
	limit := 10
	windowSize := time.Duration(500) * time.Millisecond
	limiter := NewRateLimiter(store, clock, limit, windowSize)

	result, err := limiter.Allow(t.Context(), "foo")

	if err != nil {
		t.Fatalf("unexpcted error: %v", err)
	}
	if !result {
		t.Fatalf("Allow(\"foo\") = %v; want true", result)
	}
}

func TestAllowOverLimit(t *testing.T) {
	clock := &MockClock{currentTime: time.Date(2000, 1, 1, 0, 0, 0, 0, time.UTC)}
	store := fakeStore{
		incrFunc: func() (int, error) { return 11, nil },
	}
	limit := 10
	windowSize := time.Duration(500) * time.Millisecond
	limiter := NewRateLimiter(store, clock, limit, windowSize)

	result, err := limiter.Allow(t.Context(), "foo")

	if err != nil {
		t.Fatalf("unexpcted error: %v", err)
	}
	if result {
		t.Fatalf("Allow(\"foo\") = %v; want false", result)
	}
}
