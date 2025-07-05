package store

import (
	"context"
	"time"

	"github.com/redis/go-redis/v9"
)

type RedisStore struct {
	Client *redis.Client
}

func NewRedisStore(addr string) *RedisStore {
	rdb := redis.NewClient(&redis.Options{
		Addr: addr,
	})
	return &RedisStore{Client: rdb}
}

func (s *RedisStore) Increment(ctx context.Context, key string, window time.Duration) (int, time.Duration, error) {
	val, err := s.Client.Incr(ctx, key).Result()
	if err != nil {
		return 0, 0, err
	}
	if val == 1 {
		s.Client.Expire(ctx, key, window)
	}
	ttl, err := s.Client.TTL(ctx, key).Result()
	if err != nil {
		return int(val), 0, err
	}
	return int(val), ttl, nil
}

func (s *RedisStore) Reset(ctx context.Context, key string) error {
	return s.Client.Del(ctx, key).Err()
}

func (s *RedisStore) Close() error {
	return s.Client.Close()
}

func (s *RedisStore) Ping(ctx context.Context) error {
	return s.Client.Ping(ctx).Err()
}
