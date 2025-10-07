package db

import (
	"context"
	"github.com/go-redis/redis/v8"
	"log"
	"time"
)

type NamespacedRedis struct {
	client *redis.Client
	prefix string
}

func (nr *NamespacedRedis) Set(ctx context.Context, key, value string, expirationTime time.Duration) error {
	return nr.client.Set(ctx, nr.prefix+":"+key, value, expirationTime).Err()
}

func (nr *NamespacedRedis) Get(ctx context.Context, key string) (string, error) {
	val, err := nr.client.Get(ctx, nr.prefix+":"+key).Result()
	return val, err
}

func (nr *NamespacedRedis) Del(ctx context.Context, key string) error {
	return nr.client.Del(ctx, nr.prefix+":"+key).Err()
}

func NewNamespacedRedis(client *redis.Client, prefix string) *NamespacedRedis {
	return &NamespacedRedis{
		client: client,
		prefix: prefix,
	}
}

func NewRedisClient(addr string) *redis.Client {
	client := redis.NewClient(&redis.Options{
		Addr: addr,
		DB:   0,
	})

	if _, err := client.Ping(context.Background()).Result(); err != nil {
		log.Fatalf("Failed to connect to Redis: %v", err)
	}

	return client
}
