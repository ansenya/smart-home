package handlers

import (
	"github.com/gin-gonic/gin"
	"panel-api/internal/config"
)

type Router struct {
	engine *gin.Engine
	config *config.Container

	healthHandler *healthHandler
	usersHandler  *usersHandler
}

func NewRouter(container *config.Container) Router {
	engine := gin.Default()

	return Router{
		engine:        engine,
		config:        container,
		healthHandler: newHealthHandler(),
		usersHandler:  newUsersHandler(),
	}
}

func (r *Router) registerRoutes() {
	r.healthHandler.RegisterRoutes(r.engine)
	r.usersHandler.RegisterRoutes(r.engine.Group("/panel/v1/users"))
}

func (r *Router) Run() error {
	return r.engine.Run(r.config.Server.Port)
}
