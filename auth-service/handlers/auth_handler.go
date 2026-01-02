package handlers

import (
	"auth-server/models"
	"auth-server/services"
	"errors"
	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
	"net/http"
)

type authHandler struct {
	authService services.AuthService
	jwtService  services.JWTService
}

func newAuthRouter(db *gorm.DB) *authHandler {
	return &authHandler{
		authService: services.NewAuthService(db),
		jwtService:  services.NewJwtService(),
	}
}

func (h *authHandler) RegisterRoutes(rg *gin.RouterGroup) {
	rg.GET("/me", h.Me)
	rg.POST("/login", h.Login)
	rg.POST("/logout", h.Logout)
	rg.POST("/register", h.Register)
}

func (h *authHandler) Me(c *gin.Context) {
	sid, err := c.Cookie("sid")
	if err != nil {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "invalid sid"})
		return
	}
	user, err := h.authService.Me(sid)
	if err != nil {
		c.SetCookie("sid", "", 0, "/", "", false, false)
		c.JSON(http.StatusUnauthorized, gin.H{"error": "invalid sid"})
		return
	}
	c.JSON(http.StatusOK, user)
}

func (h *authHandler) Login(c *gin.Context) {
	var request models.AuthRequest
	if err := c.ShouldBindJSON(&request); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	session, err := h.authService.Login(&request)
	if err != nil {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "wrong email/password"})
		return
	}
	c.SetCookie("sid", session.ID, int(h.jwtService.GetRefreshTokenDuration().Milliseconds()), "/", "", false, false)
	c.JSON(http.StatusOK, session)
}

func (h *authHandler) Logout(c *gin.Context) {
	// todo :)
	c.SetCookie("sid", "", 0, "/", "", false, false)
}

func (h *authHandler) Register(c *gin.Context) {
	var request models.AuthRequest
	if err := c.ShouldBindJSON(&request); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "bad body"})
		return
	}

	session, err := h.authService.Register(&request)
	if err != nil {
		if errors.Is(err, services.ErrorIncorrectPassword) {
			c.JSON(http.StatusBadRequest, gin.H{"error": "incorrect password"})
		} else {
			c.JSON(http.StatusInternalServerError, gin.H{"error": err})
		}
		return
	}

	c.SetCookie("sid", session.ID, int(h.jwtService.GetRefreshTokenDuration().Milliseconds()), "/", "", false, false)
	c.JSON(http.StatusOK, session)
}
