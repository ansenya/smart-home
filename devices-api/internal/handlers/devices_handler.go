package handlers

import (
	"devices-api/internal/models"
	"devices-api/internal/services"
	"errors"
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
)

type devicesHandler struct {
	service services.DevicesService
}

func newDevicesHandler(service services.DevicesService) *devicesHandler {
	return &devicesHandler{service: service}
}

func (h *devicesHandler) RegisterRoutes(rg *gin.RouterGroup) {
	rg.GET("", h.List)
	rg.GET("/:id", h.Get)
	rg.PUT("/:id", h.Update)
	rg.DELETE("/:id", h.Delete)
	rg.POST("/:id/capabilities/:type/set", h.SetCapability)
}

func userIDOf(c *gin.Context) (uuid.UUID, bool) {
	v, ok := c.Get("user_id")
	if !ok {
		c.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{"error": "unauthorized"})
		return uuid.Nil, false
	}
	uid, ok := v.(uuid.UUID)
	if !ok {
		c.AbortWithStatusJSON(http.StatusInternalServerError, gin.H{"error": "invalid user_id"})
		return uuid.Nil, false
	}
	return uid, true
}

func parseDeviceID(c *gin.Context) (uuid.UUID, bool) {
	id, err := uuid.Parse(c.Param("id"))
	if err != nil {
		c.AbortWithStatusJSON(http.StatusBadRequest, gin.H{"error": "invalid device id"})
		return uuid.Nil, false
	}
	return id, true
}

func (h *devicesHandler) List(c *gin.Context) {
	uid, ok := userIDOf(c)
	if !ok {
		return
	}
	devices, err := h.service.List(uid)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusOK, models.DeviceListResponse{Devices: devices})
}

func (h *devicesHandler) Get(c *gin.Context) {
	uid, ok := userIDOf(c)
	if !ok {
		return
	}
	id, ok := parseDeviceID(c)
	if !ok {
		return
	}
	device, err := h.service.Get(uid, id)
	if err != nil {
		if errors.Is(err, services.ErrDeviceNotFound) {
			c.JSON(http.StatusNotFound, gin.H{"error": "device not found"})
			return
		}
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusOK, device)
}

func (h *devicesHandler) Update(c *gin.Context) {
	uid, ok := userIDOf(c)
	if !ok {
		return
	}
	id, ok := parseDeviceID(c)
	if !ok {
		return
	}
	var req models.UpdateDeviceRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	device, err := h.service.Update(uid, id, &req)
	if err != nil {
		if errors.Is(err, services.ErrDeviceNotFound) {
			c.JSON(http.StatusNotFound, gin.H{"error": "device not found"})
			return
		}
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusOK, device)
}

func (h *devicesHandler) Delete(c *gin.Context) {
	uid, ok := userIDOf(c)
	if !ok {
		return
	}
	id, ok := parseDeviceID(c)
	if !ok {
		return
	}
	if err := h.service.Delete(uid, id); err != nil {
		if errors.Is(err, services.ErrDeviceNotFound) {
			c.JSON(http.StatusNotFound, gin.H{"error": "device not found"})
			return
		}
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	c.Status(http.StatusNoContent)
}

func (h *devicesHandler) SetCapability(c *gin.Context) {
	uid, ok := userIDOf(c)
	if !ok {
		return
	}
	id, ok := parseDeviceID(c)
	if !ok {
		return
	}
	capType := c.Param("type")
	if capType == "" {
		c.JSON(http.StatusBadRequest, gin.H{"error": "capability type required"})
		return
	}
	var req models.SetCapabilityRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	if err := h.service.SetCapabilityValue(uid, id, capType, &req); err != nil {
		if errors.Is(err, services.ErrDeviceNotFound) {
			c.JSON(http.StatusNotFound, gin.H{"error": "device not found"})
			return
		}
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	c.Status(http.StatusAccepted)
}
