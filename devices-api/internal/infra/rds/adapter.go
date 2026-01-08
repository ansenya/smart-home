package rds

import (
	"context"
	"devices-api/internal/config"
	"fmt"
	"log/slog"

	"github.com/redis/go-redis/v9"
)

type Client struct {
	cfg *config.RedisConfig
	log *slog.Logger

	Client *redis.Client
}

func NewClient(cfg *config.Container) *Client {
	cfg.Log.Info(cfg.RedisConfig.URL)
	return &Client{
		cfg: cfg.RedisConfig,
		log: cfg.Log,
	}
}

func (c *Client) Connect(ctx context.Context) error {
	ctx, cancel := context.WithTimeout(ctx, c.cfg.Timeout)
	defer cancel()

	opts := &redis.Options{
		Addr:     c.cfg.URL,
		DB:       c.cfg.DB,
		Password: c.cfg.Password,
	}

	client := redis.NewClient(opts)

	if err := client.Ping(ctx).Err(); err != nil {
		return fmt.Errorf("failed to ping redis: %w", err)
	}

	c.Client = client

	return nil
}

func (c *Client) Close() error {
	if c.Client == nil {
		return nil
	}
	return c.Client.Close()
}
