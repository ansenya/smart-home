package handlers

import (
	"github.com/gin-gonic/gin"
	"github.com/go-redis/redis/v8"
	"gorm.io/gorm"
)

type Router struct {
	healthHandler *healthHandler
	authHandler   *authHandler
	oauthHandler  *oauthHandler
}

func NewRouter(
	db *gorm.DB,
	redis *redis.Client,
) (*Router, error) {

	authHandler, err := newAuthRouter(db, redis)
	if err != nil {
		return nil, err
	}
	oauthHandler, err := newOAuthHandler(db, redis)
	if err != nil {
		return nil, err
	}
	return &Router{
		healthHandler: newHealthHandler(),
		authHandler:   authHandler,
		oauthHandler:  oauthHandler,
	}, nil
}

func (r *Router) RegisterRoutes(engine *gin.Engine) {
	// health
	r.healthHandler.RegisterRoutes(engine)

	// auth
	authGroup := engine.Group("/auth")
	r.authHandler.RegisterRoutes(authGroup)
	r.oauthHandler.RegisterRoutes(authGroup)
}
