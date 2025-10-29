package handlers

import (
	"net/http"

	"github.com/gin-gonic/gin"
)

type healthHandler struct {
}

func newHealthHandler() *healthHandler {
	return &healthHandler{}
}

func (h *healthHandler) RegisterRoutes(g *gin.Engine) {
	g.GET("/ping", func(c *gin.Context) {
		c.JSON(http.StatusOK, "pong")
	})
}
