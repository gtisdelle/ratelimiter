package ratelimiter

type Store interface {
	Increment(key string, ttlms int) (int, error)
}
