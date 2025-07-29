package ratelimiter

import (
	"sync"
	"time"
)

type memoryStore struct {
	mu    sync.Mutex
	clock Clock
	store map[string]entry
}

type entry struct {
	count     int
	expiresAt time.Time
}

func NewMemoryStore(clock Clock) Store {
	return &memoryStore{
		store: make(map[string]entry),
		clock: clock,
	}
}

func (s *memoryStore) Get(key string) (int, bool, error) {
	s.mu.Lock()
	defer s.mu.Unlock()

	entry, ok := s.store[key]
	if !ok || s.clock.Now().After(entry.expiresAt) {
		return 0, false, nil
	}

	return entry.count, true, nil
}

func (s *memoryStore) Set(key string, count int, ttlms int) error {
	s.mu.Lock()
	defer s.mu.Unlock()

	s.store[key] = entry{
		count:     count,
		expiresAt: s.clock.Now().Add(time.Duration(ttlms) * time.Millisecond),
	}

	return nil
}

func (s *memoryStore) Increment(key string, ttlms int) (int, error) {
	s.mu.Lock()
	defer s.mu.Unlock()

	e, ok := s.store[key]
	now := s.clock.Now()
	if !ok || now.After(e.expiresAt) {
		s.store[key] = entry{
			count:     1,
			expiresAt: now.Add(time.Duration(ttlms) * time.Millisecond),
		}
	} else {
		e.count++
		s.store[key] = e
	}

	return s.store[key].count, nil
}
