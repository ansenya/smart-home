package config

import (
	"log/slog"
	"os"
	"panel-api/internal/utils"
	"time"

	"github.com/joho/godotenv"
)

type Container struct {
	Server         *Server
	PostgresConfig *PostgresConfig
	Services       *Services
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
		Log: logger,
	}
}
