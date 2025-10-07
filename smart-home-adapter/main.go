package main

import (
	"devices-api/handlers"
	"devices-api/repository"
	"devices-api/services"
	"devices-api/storage"
	"devices-api/utils"
	"fmt"
	"log"

	"github.com/gin-gonic/gin"
)

var (
	Port        = fmt.Sprintf(":%s", utils.GetEnv("PORT", "8080"))
	PostgresUrl = utils.GetEnv("POSTGRES_URL", "8080")
)

func main() {
	database, err := storage.ConnectPostgres(PostgresUrl)
	if err != nil {
		log.Fatalf("failed to connect to postgres: %v", err)
	}

	// repositories
	devicesRepository := repository.NewDevicesRepo(database)

	// services
	devicesService := services.NewDevicesService(devicesRepository)

	engine := gin.Default()

	router := handlers.NewRouter(devicesService)
	router.RegisterRoutes(engine)

	if err := engine.Run(Port); err != nil {
		log.Fatalf("failed to start server: %v", err)
	}
}
