package config

import (
	"github.com/joho/godotenv"
	"log/slog"
	"panel-api/internal/utils"
)

type Container struct {
	Server Server
	Log    *slog.Logger
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
		Server: Server{
			Port: ":" + utils.GetEnv("PORT", "8080"),
		},
		Log: logger,
	}
}
