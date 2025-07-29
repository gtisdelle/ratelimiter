package ratelimiter

type Store interface {
	Get(key string) (int, bool, error)
	Set(key string, count int, ttlms int) error
	Increment(key string, ttlms int) (int, error)
}
