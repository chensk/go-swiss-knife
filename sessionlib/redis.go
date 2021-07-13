package sessionlib

import (
	redis "code.byted.org/kv/goredis/v5"
	"context"
	"errors"
	"time"
)

type RedisSessionStore struct {
	client *redis.Client
}

func NewRedisSessionStore(options *sessionOptions) (SessionStore, error) {
	if options.redisPsm == "" {
		return nil, errors.New("redis psm not found")
	}
	opt := redis.NewOption()
	opt.ReadTimeout = options.redisTimeout
	opt.WriteTimeout = options.redisTimeout
	opt.DialTimeout = options.redisTimeout
	cli, err := redis.NewClientWithOption(options.redisPsm, opt)
	if err != nil {
		return nil, err
	}

	return &RedisSessionStore{
		client: cli,
	}, nil
}

func (r *RedisSessionStore) Get(ctx context.Context, key string) (string, error) {
	v, err := r.client.Get(key).Result()
	return v, err
}

func (r *RedisSessionStore) Set(ctx context.Context, key string, value interface{}, expiration time.Duration) error {
	return r.client.Set(key, value, expiration).Err()
}

func (r *RedisSessionStore) Delete(ctx context.Context, key string) error {
	return r.client.Del(key).Err()
}
