package main

import (
	"context"
	"devices-api/internal/config"
	"devices-api/internal/handlers"
	"devices-api/internal/infra/db"
	"devices-api/internal/infra/mqtt"
	"devices-api/internal/infra/rds"
	"devices-api/internal/repositories"
	"devices-api/internal/services"
	"fmt"
	"log/slog"
	"os"
)

func main() {
	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	// config
	container := config.NewConfig()
	container.Log.Info("devices-api starting")

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

	mqttClient, err := mqtt.NewClient(container.Log)
	if err != nil {
		container.Log.Error("mqtt connect failed", slog.Any("err", err))
		os.Exit(1)
	}
	defer mqttClient.Close()

	// repositories
	sessionRepository := repositories.NewSessionRepository(dbClient.DB)
	pairingRepository := repositories.NewPairingRepository(dbClient.DB)
	pairingCache := repositories.NewPairingCache(redisClient.NewNamespacedRedis("pairing-cache"))
	devicesRepository := repositories.NewDevicesRepository(dbClient.DB)

	// services
	pairingService := services.NewPairingService(pairingRepository, pairingCache)
	devicesService := services.NewDevicesService(devicesRepository, mqttClient.Client, container.Log)
	streamService := services.NewStreamService(mqttClient.Client, container.Log)

	router := handlers.NewRouter(container, sessionRepository, pairingService, devicesService, streamService)
	if err := router.Run(); err != nil {
		container.Log.Error(fmt.Sprintf("failed to start: %s", err))
	}
}
