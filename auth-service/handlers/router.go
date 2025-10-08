package handlers

import (
	"github.com/gin-gonic/gin"
)

type Router struct {
	healthHandler HandlerInterface
}

func NewRouter() *Router {
	return &Router{
		healthHandler: newHealthHandler(),
	}
}

func (r *Router) RegisterRoutes(engine *gin.Engine) {
	r.healthHandler.RegisterRoutes(engine.Group(""))
}
