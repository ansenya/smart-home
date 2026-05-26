package handlers

import (
	"auth-server/services"
	"auth-server/utils"
	"errors"
	"log"
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/go-redis/redis/v8"
	"gorm.io/gorm"
)

type passwordResetHandler struct {
	service services.PasswordResetService
	baseURL string
}

func newPasswordResetHandler(db *gorm.DB, redisClient *redis.Client) *passwordResetHandler {
	return &passwordResetHandler{
		service: services.NewPasswordResetService(db, redisClient),
		baseURL: utils.GetEnv("AUTH_WEB_URL", "https://id.smarthome.hipahopa.ru"),
	}
}

func (h *passwordResetHandler) RegisterRoutes(rg *gin.RouterGroup) {
	rg.POST("/password-reset/request", h.Request)
	rg.POST("/password-reset/confirm", h.Confirm)
}

func (h *passwordResetHandler) Request(c *gin.Context) {
	var req struct {
		Email string `json:"email"`
	}
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "bad body"})
		return
	}

	if err := h.service.Request(req.Email, h.baseURL); err != nil {
		log.Printf("password reset request: %v", err)
	}

	// Always 200 to avoid email enumeration.
	c.JSON(http.StatusOK, gin.H{"status": "ok"})
}

func (h *passwordResetHandler) Confirm(c *gin.Context) {
	var req struct {
		Token       string `json:"token"`
		NewPassword string `json:"new_password"`
	}
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "bad body"})
		return
	}

	if err := h.service.Confirm(req.Token, req.NewPassword); err != nil {
		switch {
		case errors.Is(err, services.ErrorResetTokenInvalid):
			c.JSON(http.StatusBadRequest, gin.H{"error": "invalid or expired reset link"})
		case errors.Is(err, services.ErrorResetWeakPassword):
			c.JSON(http.StatusBadRequest, gin.H{"error": "password too weak"})
		default:
			log.Printf("password reset confirm: %v", err)
			c.JSON(http.StatusInternalServerError, gin.H{"error": "internal error"})
		}
		return
	}

	c.JSON(http.StatusOK, gin.H{"status": "ok"})
}
