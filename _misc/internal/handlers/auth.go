package handlers

import "github.com/gin-gonic/gin"

type AuthHandler struct {
}

func (h *AuthHandler) RegisterRoutes(rg *gin.RouterGroup) {
	rg.GET("/authorize", h.Authorize)
	rg.POST("/login", h.Login)
	rg.POST("/token", h.Token)
	rg.POST("/refresh", h.Refresh)
}

func (h *AuthHandler) Authorize(c *gin.Context) {}
func (h *AuthHandler) Login(c *gin.Context)     {}
func (h *AuthHandler) Token(c *gin.Context)     {}
func (h *AuthHandler) Refresh(c *gin.Context)   {}
