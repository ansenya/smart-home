package handlers

import (
	"auth-server/services"
	"github.com/gin-gonic/gin"
)

type Router struct {
	healthHandler HandlerInterface
	authHandler   HandlerInterface
}

func NewRouter(userService services.UserService,
	oauthClientsRepository services.OauthClientsService,
	oauthCodesService services.TemporaryCodeService,
	jwtService services.JWTService) *Router {
	return &Router{
		healthHandler: newHealthHandler(),
		authHandler:   newAuthRouter(userService, oauthClientsRepository, oauthCodesService, jwtService),
	}
}

func (r *Router) RegisterRoutes(engine *gin.Engine) {
	r.healthHandler.RegisterRoutes(engine.Group(""))

	authGroup := engine.Group("/auth")
	r.authHandler.RegisterRoutes(authGroup)
}
