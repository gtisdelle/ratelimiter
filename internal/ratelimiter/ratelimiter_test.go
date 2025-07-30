package ratelimiter

import (
	"testing"
	"time"
)

type fakeStore struct {
	getFunc  func() (int, bool, error)
	incrFunc func() (int, error)
}

func (s fakeStore) Get(key string) (int, bool, error) {
	return s.getFunc()
}

func (s fakeStore) Set(key string, count int, ttlms int) error {
	return nil
}

func (s fakeStore) Increment(key string, ttlms int) (int, error) {
	return s.incrFunc()
}

func TestAllowUnderLimit(t *testing.T) {
	clock := &MockClock{currentTime: time.Date(2000, 1, 1, 0, 0, 0, 0, time.UTC)}
	store := fakeStore{
		getFunc:  func() (int, bool, error) { return 1, true, nil },
		incrFunc: func() (int, error) { return 1, nil },
	}
	limit := 10
	windowSize := time.Duration(500) * time.Millisecond
	limiter := NewRateLimiter(store, clock, limit, windowSize)

	result, err := limiter.Allow("foo")

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
		getFunc:  func() (int, bool, error) { return 10, true, nil },
		incrFunc: func() (int, error) { return 11, nil },
	}
	limit := 10
	windowSize := time.Duration(500) * time.Millisecond
	limiter := NewRateLimiter(store, clock, limit, windowSize)

	result, err := limiter.Allow("foo")

	if err != nil {
		t.Fatalf("unexpcted error: %v", err)
	}
	if result {
		t.Fatalf("Allow(\"foo\") = %v; want false", result)
	}
}
