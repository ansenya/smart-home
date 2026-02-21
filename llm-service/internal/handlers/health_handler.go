package handlers

import (
	"net/http"

	"github.com/gin-gonic/gin"
)

type healthHandler struct {
}

func (h *healthHandler) RegisterRoutes(rg *gin.RouterGroup) {
	rg.GET("/ping", h.HandlePing)
}

func (h *healthHandler) HandlePing(c *gin.Context) {
	c.JSON(http.StatusOK, "pong")
}

func NewHealthHandler() Handler {
	return &healthHandler{}
}
