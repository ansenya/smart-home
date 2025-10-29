package main

import (
	"auth-server/handlers"
	"auth-server/repository"
	"auth-server/services"
	"auth-server/storage"
	"auth-server/utils"
	"fmt"
	"github.com/gin-contrib/cors"
	"github.com/gin-gonic/gin"
	"log"
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

	// repositories
	userRepository := repository.NewUserRepository(database)

	// services
	userService := services.NewUserService(userRepository)
	oauthClientsService := services.NewOauthClientsService(database, redisClient)
	oauthCodeService := services.NewOauthCodeService(redisClient)
	jwtService, err := services.NewJwtService()
	if err != nil {
		log.Fatalf("failed to create jwt service: %v", err)
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

	router, err := handlers.NewRouter(database, userService, oauthClientsService, oauthCodeService, jwtService)
	if err != nil {
		log.Fatalf("failed to create router: %s", err)
	}
	router.RegisterRoutes(engine)

	if err := engine.Run(Port); err != nil {
		log.Fatalf("failed to start server: %v", err)
	}
}
