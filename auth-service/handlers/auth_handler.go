package handlers

import (
	"auth-server/models"
	"auth-server/services"
	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
	"net/http"
)

type authHandler struct {
	authService services.AuthService
	jwtService  services.JWTService
}

func newAuthRouter(db *gorm.DB) (*authHandler, error) {
	authService, err := services.NewAuthService(db)
	if err != nil {
		return nil, err
	}
	jwtService, err := services.NewJwtService()
	if err != nil {
		return nil, err
	}
	return &authHandler{
		authService: authService,
		jwtService:  jwtService,
	}, nil
}

func (h *authHandler) RegisterRoutes(rg *gin.RouterGroup) {
	rg.POST("/login", h.Login)
	rg.POST("/register", h.Register)
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

func (h *authHandler) Register(c *gin.Context) {
	var request models.AuthRequest
	if err := c.ShouldBindJSON(&request); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "bad body"})
		return
	}

	user, err := h.authService.Register(&request)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err})
		return
	}

	c.JSON(http.StatusCreated, user)
}
