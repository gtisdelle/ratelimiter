package limit

import (
	"context"
	"fmt"
	"time"

	"github.com/redis/go-redis/v9"
)

type clock interface {
	now() time.Time
}

type Config struct {
	BucketSize int
	Rate       int
}

type redisStore struct {
	rdb   *redis.Client
	clock clock
	cfg   Config
}

func NewRedisStore(rdb *redis.Client, clock clock, cfg Config) store {
	return &redisStore{
		rdb:   rdb,
		clock: clock,
		cfg:   cfg,
	}
}

func (s *redisStore) Allow(ctx context.Context, key string, hits uint64) (bool, int, error) {
	lua := redis.NewScript(`
	local key = KEYS[1]
	local rate = tonumber(ARGV[1])
	local capacity = tonumber(ARGV[2])
	local now = tonumber(ARGV[3])
	local hits = tonumber(ARGV[4])
	local result = redis.call("HMGET", key, "tokens", "lastRefill")
	local tokens = tonumber(result[1])
	local lastRefill = tonumber(result[2])

	-- all data nil is interpreted to mean no value for tenant is set yet
	if tokens == nil and lastRefill == nil then
	  redis.log(redis.LOG_WARNING, "Miss - tokens: " .. tostring(tokens) .. ", lastRefill: " .. tostring(lastRefill))
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
	if tokens >= hits then
	  allow = 1
	else
	  allow = 0
	end
	tokens = math.max(tokens - hits, 0)
	lastRefill = now

	redis.call("HMSET", key, "tokens", tokens, "lastRefill", lastRefill)
	local ttl = math.ceil(capacity / rate)
	redis.call("EXPIRE", key, ttl)
	
	-- redis.log(redis.LOG_WARNING, "Usage consumed - tokens: " .. tostring(tokens) .. ", lastRefill: " .. tostring(lastRefill) .. ", allow: " .. allow)

	return { allow, tokens }
	`)

	keys := []string{key}
	args := []any{s.cfg.Rate, s.cfg.BucketSize, s.clock.now().UnixMilli(), hits}
	result, err := lua.Eval(ctx, s.rdb, keys, args).Result()
	if err != nil {
		return false, 0, fmt.Errorf("token bucket lua script: %w", err)
	}
	arr, ok := result.([]any)
	if !ok {
		return false, 0, fmt.Errorf("invalid response type from lua script: %T", result)
	}
	if len(arr) != 2 {
		return false, 0, fmt.Errorf("invalid response size from lua script: %d", len(arr))
	}
	allow, ok := arr[0].(int64)
	if !ok {
		return false, 0, fmt.Errorf("invalid first array value type from lua script: %T", arr[0])
	}
	remaining, ok := arr[1].(int64)
	if !ok {
		return false, 0, fmt.Errorf("invalid first array value type from lua script: %T", arr[1])
	}

	return allow > 0, int(remaining), nil
}
