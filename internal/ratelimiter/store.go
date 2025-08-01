package ratelimiter

type Store interface {
	Increment(key string, ttlms int) (int, error)

	// TODO: remove these unused functions
	Get(key string) (int, bool, error)
	Set(key string, count int, ttlms int) error
}
