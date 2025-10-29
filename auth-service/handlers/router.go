package handlers

import (
	"auth-server/services"
	"github.com/gin-gonic/gin"
)

type Router struct {
	healthHandler *healthHandler
	authHandler   *authHandler
}

func NewRouter(
	userService services.UserService,
	oauthClientsRepository services.OauthClientsService,
	oauthCodesService services.TemporaryCodeService,
	jwtService services.JWTService,
) *Router {

	return &Router{
		healthHandler: newHealthHandler(),
		authHandler:   newAuthRouter(userService, oauthClientsRepository, oauthCodesService, jwtService),
	}
}

func (r *Router) RegisterRoutes(engine *gin.Engine) {
	// health
	r.healthHandler.RegisterRoutes(engine)

	// auth
	authGroup := engine.Group("/auth")
	r.authHandler.RegisterRoutes(authGroup)
}
