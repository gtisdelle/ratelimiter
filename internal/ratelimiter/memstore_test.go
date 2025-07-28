package ratelimiter

import "testing"

func TestGet(t *testing.T) {
	s := NewMemoryStore(NewClock())
	s.Set("foo", 3, 60)

	result, err := s.Get("foo")

	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if result != 3 {
		t.Fatalf("Get(foo) = %d; want 3", result)
	}
}
