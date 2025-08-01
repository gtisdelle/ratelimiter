package ratelimiter

import (
	"testing"
	"time"
)

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

func TestDoubleIncrement(t *testing.T) {
	clock := &MockClock{
		currentTime: time.Date(2000, 1, 1, 0, 0, 0, 0, time.UTC),
	}
	s := NewMemoryStore(clock)

	_, err := s.Increment("foo", 60000)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	result, err := s.Increment("foo", 60000)

	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if result != 2 {
		t.Fatalf("Increment(foo, 60)<2x> = %d; want 2", result)
	}
}

func TestDoubleIncrementTTL(t *testing.T) {
	clock := &MockClock{
		currentTime: time.Date(2000, 1, 1, 0, 0, 0, 0, time.UTC),
	}
	s := NewMemoryStore(clock)

	_, err := s.Increment("foo", 60000)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	clock.Advance(time.Duration(61) * time.Second)
	result, err := s.Increment("foo", 60000)

	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if result != 1 {
		t.Fatalf("Increment(foo, 60)<2x> = %d; want 1", result)
	}
}
