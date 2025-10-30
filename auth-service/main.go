package main

import (
	"auth-server/handlers"
	"auth-server/storage"
	"auth-server/utils"
	"fmt"
	"github.com/gin-contrib/cors"
	"github.com/gin-gonic/gin"
	"log"
)

var (
	Port        = fmt.Sprintf(":%s", utils.GetEnv("PORT", "8080"))
	PostgresUrl = utils.GetEnv("POSTGRES_URL", "host=localhost user=user password=password dbname=smart-home port=5432 sslmode=disable")
	RedisUrl    = utils.GetEnv("REDIS_URL", "localhost:6379")
)

func main() {
	database, err := storage.ConnectPostgres(PostgresUrl)
	if err != nil {
		log.Fatalf("failed to connect to postgres: %v", err)
	}

	redisClient, err := storage.NewRedisClient(RedisUrl)
	if err != nil {
		log.Fatalf("failed to connect to redis: %v", err)
	}

	engine := gin.Default()
	engine.Use(cors.New(cors.Config{
		AllowOrigins:     []string{"https://smarthome.hipahopa.ru", "http://localhost:5173"},
		AllowMethods:     []string{"GET", "POST", "OPTIONS"},
		AllowHeaders:     []string{"Origin", "Content-Type", "Accept", "Authorization"},
		ExposeHeaders:    []string{"Content-Length"},
		AllowCredentials: true,
		MaxAge:           86400,
	}))

	router := handlers.NewRouter(database, redisClient)
	router.RegisterRoutes(engine)

	if err := engine.Run(Port); err != nil {
		log.Fatalf("failed to start server: %v", err)
	}
}
