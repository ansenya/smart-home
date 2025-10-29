package handlers

import (
	"auth-server/services"
	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
)

type Router struct {
	healthHandler *healthHandler
	authHandler   *authHandler
}

func NewRouter(
	db *gorm.DB,
	userService services.UserService,
	oauthClientsRepository services.OauthService,
	oauthCodesService services.TemporaryCodeService,
	jwtService services.JWTService,
) (*Router, error) {

	authHandler, err := newAuthRouter(db, userService, oauthClientsRepository, oauthCodesService, jwtService)
	if err != nil {
		return nil, err
	}
	return &Router{
		healthHandler: newHealthHandler(),
		authHandler:   authHandler,
	}, nil
}

func (r *Router) RegisterRoutes(engine *gin.Engine) {
	// health
	r.healthHandler.RegisterRoutes(engine)

	// auth
	authGroup := engine.Group("/auth")
	r.authHandler.RegisterRoutes(authGroup)
}
