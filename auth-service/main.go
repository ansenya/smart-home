package main

import (
	"auth-server/handlers"
	"auth-server/service"
	"auth-server/storage"
	"auth-server/utils"
	"fmt"
	"github.com/gin-contrib/cors"
	"github.com/gin-gonic/gin"
	"log"
	"os"
	"time"
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
		AllowAllOrigins:  true,
		AllowMethods:     []string{"*"},
		AllowHeaders:     []string{"*"},
		ExposeHeaders:    []string{"*"},
		AllowCredentials: true,
		MaxAge:           12 * time.Hour,
	}))

	// init smtp
	smtpConfig := service.SmtpConfig{
		Host:     os.Getenv("SMTP_HOST"),
		Port:     os.Getenv("SMTP_PORT"),
		Password: os.Getenv("SMTP_PASSWORD"),
		Username: os.Getenv("SMTP_USERNAME"),
	}

	// register routes
	handlers.RegisterAuthRoutes(engine, database, redisClient, smtpConfig)

	if err := engine.Run(Port); err != nil {
		log.Fatalf("failed to start server: %v", err)
	}
}
