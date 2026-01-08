package main

import (
	"context"
	"fmt"
	"log/slog"
	"os"
	"panel-api/internal/config"
	"panel-api/internal/handlers"
	"panel-api/internal/infra/db"
	"panel-api/internal/repositories"
	"panel-api/internal/services"
	"time"
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

	// repositories
	sessionRepository := repositories.NewSessionRepository(dbClient.DB)
	usersRepository := repositories.NewUsersRepository(dbClient.DB)

	// services
	oauthService := services.NewOauthService(services.OauthConfig{
		BaseURL:       "https://api.id.smarthome.hipahopa.ru",
		TokenEndpoint: "/oauth/token",
		UserEndpoint:  "/oauth/userinfo",
		Timeout:       5 * time.Second,
	}, nil)
	usersService := services.NewUsersService(container.Log, sessionRepository, usersRepository)

	container.Services = &config.Services{
		OauthService: oauthService,
		UsersService: usersService,
	}

	router := handlers.NewRouter(container)
	if err := router.Run(); err != nil {
		container.Log.Error(fmt.Sprintf("failed to start: %s", err))
	}
}
