package handlers

import (
	"github.com/gin-gonic/gin"
	"registry/config"
)

type Router struct {
	engine        *gin.Engine
	config        *config.Container
	healthHandler *healthHandler
}

func (r *Router) Run() error {
	return r.engine.Run(r.config.Server.Port)
}

func (r *Router) registerRoutes() {
	r.healthHandler.RegisterRoutes(r.engine)
}

func NewRouter(config *config.Container) *Router {
	router := Router{
		engine:        gin.Default(),
		config:        config,
		healthHandler: newHealthHandler(),
	}
	router.registerRoutes()
	return &router
}
