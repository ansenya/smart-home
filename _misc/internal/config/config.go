package config

import (
	"os"
	"time"
)

type Config struct {
	Server struct {
		Port string
	}
	Postgres struct {
		DSN string
	}
	Redis struct {
		Addr string
	}
	Tokens struct {
		AccessPrivateKeyPath  string
		AccessPublicKeyPath   string
		RefreshPrivateKeyPath string
		RefreshPublicKeyPath  string
		AccessTTL             time.Duration
		RefreshTTL            time.Duration
	}
}

func Load() (*Config, error) {
	cfg := &Config{}

	cfg.Server.Port = ":" + getEnv("PORT", "8080")
	cfg.Postgres.DSN = getEnv("POSTGRES_DSN", "")
	cfg.Redis.Addr = getEnv("REDIS_ADDR", "localhost:6379")

	cfg.Tokens.AccessPrivateKeyPath = getEnv("ACCESS_PRIVATE_KEY", "/secrets/access_private.pem")
	cfg.Tokens.AccessPublicKeyPath = getEnv("ACCESS_PUBLIC_KEY", "/secrets/access_public.pem")
	cfg.Tokens.RefreshPrivateKeyPath = getEnv("REFRESH_PRIVATE_KEY", "/secrets/refresh_private.pem")
	cfg.Tokens.RefreshPublicKeyPath = getEnv("REFRESH_PUBLIC_KEY", "/secrets/refresh_public.pem")

	cfg.Tokens.AccessTTL = 15 * time.Minute
	cfg.Tokens.RefreshTTL = 30 * 24 * time.Hour

	return cfg, nil
}

func getEnv(key, fallback string) string {
	if v := os.Getenv(key); v != "" {
		return v
	}
	return fallback
}
