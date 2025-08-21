package ratelimiter

import (
	"context"
	"fmt"

	"github.com/redis/go-redis/v9"
)

type Config struct {
	BucketSize int
	Rate       int
}

type redisStore struct {
	rdb   *redis.Client
	clock Clock
	cfg   Config
}

func NewRedisStore(rdb *redis.Client, clock Clock, cfg Config) Store {
	return &redisStore{
		rdb:   rdb,
		clock: clock,
		cfg:   cfg,
	}
}

func (s *redisStore) Allow(ctx context.Context, key string) (bool, error) {
	lua := redis.NewScript(`
	local key = KEYS[1]
	local rate = tonumber(ARGV[1])
	local capacity = tonumber(ARGV[2])
	local now = tonumber(ARGV[3])
	local result = redis.call("HMGET", key, "tokens", "lastRefill")
	local tokens = tonumber(result[1])
	local lastRefill = tonumber(result[2])

	-- all data nil is interpreted to mean no value for tenant is set yet
	if tokens == nil and lastRefill == nil then
	  tokens = capacity
	  lastRefill = now
	end
	-- one or the other being nil breaks the invariant of both always being set
	if tokens == nil or lastRefill == nil then
	  redis.log(redis.LOG_WARNING, "Invariant broken - tokens: " .. tostring(tokens) .. ", lastRefill: " .. tostring(lastRefill))
	  return false
	end

	local delta = math.max((now - lastRefill) / 1000, 0)
	tokens = math.min(tokens + (rate * delta), capacity)

	local allow = 0
	if tokens >= 1 then
	  allow = 1
	else
	  allow = 0
	end
	tokens = math.max(tokens - 1, 0)
	lastRefill = now

	redis.call("HMSET", key, "tokens", tokens, "lastRefill", lastRefill)
	local ttl = (capacity / rate) * 1000
	redis.call("EXPIRE", key, ttl)

	return allow
	`)

	keys := []string{key}
	args := []any{s.cfg.Rate, s.cfg.BucketSize, s.clock.Now().UnixMilli()}
	result, err := lua.Eval(ctx, s.rdb, keys, args).Result()
	if err != nil {
		return false, fmt.Errorf("token bucket lua script: %w", err)
	}
	allow, ok := result.(int64)
	if !ok {
		return false, fmt.Errorf("invalid response type from lua script: %T", result)
	}
	return allow > 0, nil
}

var _ Store = &redisStore{}
