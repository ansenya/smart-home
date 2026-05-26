package handlers

import (
	"auth-server/models"
	"auth-server/services"
	"errors"
	"log"
	"net/http"
	"strings"

	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
)

func setSessionCookie(c *gin.Context, value string, maxAgeSec int) {
	secure := !strings.HasSuffix(DomainName, ".local")
	c.SetSameSite(http.SameSiteLaxMode)
	c.SetCookie(SessionIDName, value, maxAgeSec, "/", DomainName, secure, true)
}

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
	sid, err := c.Cookie(SessionIDName)
	if err != nil {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "invalid sid"})
		return
	}
	user, err := h.authService.Me(sid)
	if err != nil {
		setSessionCookie(c, "", -1)
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
		if errors.Is(err, services.ErrorInvalidCredentials) {
			c.JSON(http.StatusUnauthorized, gin.H{"error": "invalid email or password"})
			return
		}
		log.Printf("login failed: %v", err)
		c.JSON(http.StatusInternalServerError, gin.H{"error": "internal error"})
		return
	}
	setSessionCookie(c, session.ID, int(h.jwtService.GetRefreshTokenDuration().Seconds()))
	c.JSON(http.StatusOK, session)
}

func (h *authHandler) Logout(c *gin.Context) {
	if sid, err := c.Cookie(SessionIDName); err == nil && sid != "" {
		if err := h.authService.Logout(sid); err != nil {
			log.Printf("logout: failed to invalidate session: %v", err)
		}
	}
	setSessionCookie(c, "", -1)
	c.Status(http.StatusNoContent)
}

func (h *authHandler) Register(c *gin.Context) {
	var request models.AuthRequest
	if err := c.ShouldBindJSON(&request); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "bad body"})
		return
	}

	session, err := h.authService.Register(&request)
	if err != nil {
		switch {
		case errors.Is(err, services.ErrorEmailExists):
			c.JSON(http.StatusConflict, gin.H{"error": "email already registered"})
		case errors.Is(err, services.ErrorInvalidEmail):
			c.JSON(http.StatusBadRequest, gin.H{"error": "invalid email"})
		case errors.Is(err, services.ErrorIncorrectPassword):
			c.JSON(http.StatusBadRequest, gin.H{"error": "incorrect password"})
		default:
			log.Printf("register failed: %v", err)
			c.JSON(http.StatusInternalServerError, gin.H{"error": "internal error"})
		}
		return
	}

	setSessionCookie(c, session.ID, int(h.jwtService.GetRefreshTokenDuration().Seconds()))
	c.JSON(http.StatusOK, session)
}
