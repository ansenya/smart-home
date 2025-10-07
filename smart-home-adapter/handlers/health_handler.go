package handlers

import (
	"net/http"

	"github.com/gin-gonic/gin"
)

type healthHandler struct {
}

func newHealthHandler() HandlerInterface {
	return &healthHandler{}
}

func (h *healthHandler) RegisterRoutes(g *gin.RouterGroup) {
	g.GET("/ping", func(c *gin.Context) {
		c.JSON(http.StatusOK, "pong")
	})
}
