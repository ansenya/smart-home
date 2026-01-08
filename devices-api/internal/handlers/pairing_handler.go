package handlers

import (
	"devices-api/internal/services"
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
)

type pairingHandler struct {
	service services.PairingService
}

func newPairingHandler(service services.PairingService) *pairingHandler {
	return &pairingHandler{
		service: service,
	}
}

func (h *pairingHandler) RegisterRoutes(usersGroup *gin.RouterGroup, devicesGroup *gin.RouterGroup) {
	usersGroup.POST("/start", h.Start)
	devicesGroup.POST("/confirm", h.Confirm)
}

func (h *pairingHandler) Start(c *gin.Context) {
	userID := c.MustGet("user_id").(uuid.UUID)

	code, expires, err := h.service.StartPairing(userID)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"code":       code,
		"expires_in": expires,
	})
}
func (h *pairingHandler) Confirm(c *gin.Context) {
	code := c.Query("code")
	deviceUID := c.Query("device_uid")

	if code == "" || deviceUID == "" {
		c.JSON(http.StatusBadRequest, gin.H{"error": "code and device_uid required"})
		return
	}

	if err := h.service.ConfirmPairing(code, deviceUID); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{"status": "ok"})
}
