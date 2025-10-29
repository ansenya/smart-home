package handlers

import (
	"github.com/gin-gonic/gin"
	"net/http"
)

type HealthHandler struct {
}

func (h *HealthHandler) RegisterRoutes(engine *gin.Engine) {
	engine.GET("/ping", func(c *gin.Context) {
		c.JSON(http.StatusOK, "pong")
	})
}
