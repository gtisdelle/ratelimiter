package ratelimiter

import (
	"context"
	"time"

	"github.com/redis/go-redis/v9"
)

type redisStore struct {
	rdb *redis.Client
}

func NewRedisStore(rdb *redis.Client) Store {
	return &redisStore{
		rdb: rdb,
	}
}

func (s *redisStore) Increment(key string, ttl time.Duration) (int, error) {
	ok, err := s.rdb.SetNX(context.TODO(), key, 1, ttl).Result()
	if err != nil {
		return 0, err
	}
	if ok {
		return 1, nil
	}

	val, err := s.rdb.Incr(context.TODO(), key).Result()
	if err != nil {
		return 0, err
	}

	return int(val), nil
}
