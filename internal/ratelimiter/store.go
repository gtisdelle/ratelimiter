package ratelimiter

type Store interface {
	Get(key string) (int, error)
	Set(key string, count int, ttlseconds int) error
	Increment(key string) (int, error)
}
