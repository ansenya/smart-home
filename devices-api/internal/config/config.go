package config

import (
	"device-service/internal/utils"
	"log/slog"
	"os"
	"time"

	"github.com/joho/godotenv"
)

type Container struct {
	Server         *Server
	PostgresConfig *PostgresConfig
	RedisConfig    *RedisConfig
	Log            *slog.Logger
}

type Server struct {
	Port string
}

func NewConfig() *Container {
	logger := initLogger()

	err := godotenv.Load()
	if err != nil {
		logger.Warn(".env file not found")
	}

	return &Container{
		Server: &Server{
			Port: ":" + utils.GetEnv("PORT", "8080"),
		},
		PostgresConfig: &PostgresConfig{
			URL:              os.Getenv("POSTGRES_URL"),
			MaxOpenConns:     25,
			MaxIdleConns:     5,
			ConnMaxLifetime:  30 * time.Minute,
			StatementTimeout: 5 * time.Second,
			LockTimeout:      1 * time.Second,
		},
		RedisConfig: &RedisConfig{
			URL:      utils.GetEnv("REDIS_URL", "redis:6379"),
			Password: utils.GetEnv("REDIS_PASSWORD", ""),
			DB:       0,
			Timeout:  5 * time.Second,
		},
		Log: logger,
	}
}
