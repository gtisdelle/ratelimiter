package ratelimiter

import (
	"testing"
	"time"
)

func TestSetThenGet(t *testing.T) {
	clock := &MockClock{
		currentTime: time.Now(),
	}
	s := NewMemoryStore(clock)
	s.Set("foo", 3, 60000)

	result, ok, err := s.Get("foo")

	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if !ok || result != 3 {
		t.Fatalf("Get(foo) = (%d, %v); want (3, true)", result, ok)
	}
}

func TestSetThenGetTTL(t *testing.T) {
	clock := &MockClock{
		currentTime: time.Now(),
	}
	s := NewMemoryStore(clock)
	s.Set("foo", 3, 500)

	clock.Advance(time.Duration(501) * time.Millisecond)
	result, ok, err := s.Get("foo")

	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if ok || result != 0 {
		t.Fatalf("Get(foo) = (%d, %v); want (0, false)", result, ok)
	}
}

func TestIncrementNotExists(t *testing.T) {
	clock := &MockClock{
		currentTime: time.Now(),
	}
	s := NewMemoryStore(clock)

	result, err := s.Increment("foo", int(time.Duration(60)*time.Second))

	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if result != 1 {
		t.Fatalf("Increment(foo, 60) = %d; want 1", result)
	}
}

func TestSetThenIncrement(t *testing.T) {
	clock := &MockClock{
		currentTime: time.Now(),
	}
	s := NewMemoryStore(clock)
	s.Set("foo", 3, 60000)

	result, err := s.Increment("foo", 60000)

	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if result != 4 {
		t.Fatalf("Increment(foo, 60) = %d; want 4", result)
	}
}
