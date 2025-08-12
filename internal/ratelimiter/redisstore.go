package ratelimiter

import (
	"context"
	"fmt"
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

func (s *redisStore) Increment(ctx context.Context, key string, ttl time.Duration) (int, error) {
	commands, err := s.rdb.Pipelined(ctx, func(p redis.Pipeliner) error {
		p.SetNX(ctx, key, 0, ttl)
		p.Incr(ctx, key)
		return nil
	})
	if err != nil {
		return 0, fmt.Errorf("pipeline exec: %w", err)
	}
	incrCmd, ok := commands[1].(*redis.IntCmd)
	if !ok {
		return 0, fmt.Errorf("unexpected pipeline response type: %T", commands[1])
	}
	if err := incrCmd.Err(); err != nil {
		return 0, fmt.Errorf("redis INCR command: %w", err)
	}
	return int(incrCmd.Val()), nil
}
