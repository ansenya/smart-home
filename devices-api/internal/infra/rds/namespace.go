package rds

import (
	"context"
	"time"

	"github.com/redis/go-redis/v9"
)

type NamespacedRedis interface {
	Set(ctx context.Context, key string, value any, ttl time.Duration) error
	Get(ctx context.Context, key string) (string, error)
	Del(ctx context.Context, key string) error
}

type namespacedRedis struct {
	client    *redis.Client
	namespace string
}

func NewNamespacedRedis(client *redis.Client, namespace string) NamespacedRedis {
	return &namespacedRedis{
		client:    client,
		namespace: namespace,
	}
}

func (r *namespacedRedis) k(key string) string {
	return r.namespace + ":" + key
}

func (r *namespacedRedis) Set(ctx context.Context, key string, value any, ttl time.Duration) error {
	return r.client.Set(ctx, r.k(key), value, ttl).Err()
}

func (r *namespacedRedis) Get(ctx context.Context, key string) (string, error) {
	return r.client.Get(ctx, r.k(key)).Result()
}

func (r *namespacedRedis) Del(ctx context.Context, key string) error {
	return r.client.Del(ctx, r.k(key)).Err()
}
