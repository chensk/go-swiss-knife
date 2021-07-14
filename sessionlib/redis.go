package sessionlib

import (
	"context"
	"errors"
	"github.com/go-redis/redis/v8"
	"time"
)

// InMemorySessionStore stores the session in redis.
type RedisSessionStore struct {
	client        *redis.Client
	clusterClient *redis.ClusterClient
}

func NewRedisSessionStore(options *sessionOptions) (SessionStore, error) {
	if len(options.RedisClusters) == 0 {
		return nil, errors.New("redis cluster not found")
	}
	if len(options.RedisClusters) == 1 {
		cli := redis.NewClient(&redis.Options{
			Addr: options.RedisClusters[0],
		})
		return &RedisSessionStore{
			client: cli,
		}, nil
	} else {
		cli := redis.NewClusterClient(&redis.ClusterOptions{
			Addrs: options.RedisClusters,
		})
		return &RedisSessionStore{
			clusterClient: cli,
		}, nil
	}
}

func (r *RedisSessionStore) Get(ctx context.Context, key string) (string, error) {
	if r.client != nil {
		v, err := r.client.Get(ctx, key).Result()
		return v, err
	} else {
		v, err := r.clusterClient.Get(ctx, key).Result()
		return v, err
	}
}

func (r *RedisSessionStore) Set(ctx context.Context, key string, value string, expiration time.Duration) error {
	if r.client != nil {
		return r.client.Set(ctx, key, value, expiration).Err()
	} else {
		return r.clusterClient.Set(ctx, key, value, expiration).Err()
	}
}

func (r *RedisSessionStore) Delete(ctx context.Context, key string) error {
	if r.client != nil {
		return r.client.Del(ctx, key).Err()
	} else {
		return r.clusterClient.Del(ctx, key).Err()
	}
}
