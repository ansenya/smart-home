package app

import (
	"auth-server/internal/api/handlers"
	"auth-server/internal/config"
	"auth-server/internal/domain/oauth"
	"auth-server/internal/domain/token"
	"auth-server/internal/domain/user"
	"auth-server/internal/infra/postgres"
	"auth-server/internal/infra/redis"
)

type Container struct {
	Router *handlers.Router
}

func NewContainer(cfg *config.Config) (*Container, error) {
	// DB
	db, err := postgres.Connect(cfg.Postgres.DSN)
	if err != nil {
		return nil, err
	}

	redisClient, err := redis.New(cfg.Redis.Addr)
	if err != nil {
		return nil, err
	}

	// Repositories
	userRepo := postgres.NewUserRepository(db)
	clientRepo := postgres.NewOauthClientRepository(db)
	codeRepo := redis.NewOAuthCodeRepository(redisClient)

	// Services (domain)
	tokenService := token.NewService(cfg.JWT)
	userService := user.NewService(userRepo)
	oauthService := oauth.NewService(userRepo, clientRepo, codeRepo, tokenService)

	// API router
	router := handlers.NewRouter(userService, oauthService)

	return &Container{Router: router}, nil
}
