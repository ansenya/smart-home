package main

import (
	"context"
	"device-service/internal/config"
	"device-service/internal/handlers"
	"device-service/internal/infra/db"
	"device-service/internal/infra/rds"
	"fmt"
	"log/slog"
	"os"
)

func main() {
	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	// config
	container := config.NewConfig()
	container.Log.Info("notification service started")

	// DB
	dbClient := db.NewClient(container)
	if err := dbClient.Connect(ctx); err != nil {
		container.Log.Error("db connect failed", slog.Any("err", err))
		os.Exit(1)
	}
	defer dbClient.Close()

	redisClient := rds.NewClient(container)
	if err := redisClient.Connect(ctx); err != nil {
		container.Log.Error("redis connect failed", slog.Any("err", err))
		os.Exit(1)
	}
	defer redisClient.Close()

	router := handlers.NewRouter(container)
	if err := router.Run(); err != nil {
		container.Log.Error(fmt.Sprintf("failed to start: %s", err))
	}
}
