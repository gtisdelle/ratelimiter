package ratelimiter

import "time"

type Store interface {
	Increment(key string, ttl time.Duration) (int, error)
}
