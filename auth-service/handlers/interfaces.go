package handlers

import "github.com/gin-gonic/gin"

type HandlerInterface interface {
	RegisterRoutes(g *gin.RouterGroup)
}
