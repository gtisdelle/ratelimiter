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
	count      int
	experiesAt time.Time
}

func NewMemoryStore(clock Clock) Store {
	return &memoryStore{
		store: make(map[string]entry),
	}
}

func (s *memoryStore) Get(key string) (int, error) {
	return 3, nil
}

func (s *memoryStore) Set(key string, count int, ttlseconds int) error {
	return nil
}

func (s *memoryStore) Increment(key string) (int, error) {
	return 0, nil
}
