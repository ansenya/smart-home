package handlers

import (
	"github.com/gin-gonic/gin"
	"github.com/go-redis/redis/v8"
	"gorm.io/gorm"
)

// todo надо перейти на разные домены для разных кусков бэкенда

const (
	SessionIDName = "oauth-sid"
	DomainName    = "api.smarthome.hipahopa.ru"
)

type Router struct {
	healthHandler *healthHandler
	authHandler   *authHandler
	oauthHandler  *oauthHandler
}

func NewRouter(db *gorm.DB, redisClient *redis.Client) *Router {
	return &Router{
		healthHandler: newHealthHandler(),
		authHandler:   newAuthRouter(db),
		oauthHandler:  newOauthHandler(db, redisClient),
	}
}

func (r *Router) RegisterRoutes(engine *gin.Engine) {
	// health
	r.healthHandler.RegisterRoutes(engine)

	// auth
	authGroup := engine.Group("/auth")
	r.authHandler.RegisterRoutes(authGroup)

	// oauth
	oauthGroup := engine.Group("/oauth")
	r.oauthHandler.RegisterRoutes(oauthGroup)
}
