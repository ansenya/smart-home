package main

import (
	"context"
	"fmt"
	"llm-service/internal/config"
	"llm-service/internal/handlers"
	"llm-service/internal/infra/db"
	"llm-service/internal/repositories"
	"llm-service/internal/services"
	"log/slog"
	"os"
)

func main() {
	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	cfg := config.NewConfig()
	cfg.Log.Info("faq service started")

	// DB
	dbClient := db.NewClient(cfg)
	if err := dbClient.Connect(ctx); err != nil {
		cfg.Log.Error("db connect failed", slog.Any("err", err))
		os.Exit(1)
	}
	defer dbClient.Close()

	// services
	repos := repositories.NewContainer(dbClient)
	svcs := services.NewContainer(repos)

	// router
	router := handlers.NewRouter(cfg, svcs, repos)
	if err := router.Run(); err != nil {
		cfg.Log.Error(fmt.Sprintf("failed to start: %s", err))
	}

	defer dbClient.Close()
}
