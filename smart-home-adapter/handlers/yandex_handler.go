package handlers

import (
	"devices-api/models"
	"devices-api/services"
	"fmt"
	"log"
	"net/http"

	"github.com/gin-gonic/gin"
)

type response struct {
	RequestID string  `json:"request_id"`
	Payload   payload `json:"payload"`
}

type payload struct {
	UserID  string `json:"user_id,omitempty"`
	Devices []models.Device
}

type yandexHandler struct {
	devicesService services.DevicesService
}

func (h *yandexHandler) RegisterRoutes(g *gin.RouterGroup) {
	g.GET("", h.handleAliveStatus)

	protected := g.Group("/")
	protected.Use() // должен быть подтвержденный акк
	protected.POST("/user/unlink", h.handleUnlink)
	protected.GET("/user/devices", h.handleDevices)
	protected.POST("/user/devices/query", h.handleDevicesQuery)
	protected.POST("/yandex/v1.0/user/devices/action", h.handleDevicesAction)
}

func (h *yandexHandler) handleAliveStatus(c *gin.Context) {
	c.Status(http.StatusOK)
}

func (h *yandexHandler) handleUnlink(c *gin.Context) {

}

func (h *yandexHandler) handleDevices(c *gin.Context) {
	requestID := getValueFromContext(c, "requestID")
	userID := getValueFromContext(c, "userID")

	devices, err := h.devicesService.GetUserDevices(userID)
	if err != nil {
		log.Printf("cannot list devices: %s", err.Error())
		c.AbortWithStatusJSON(http.StatusInternalServerError, gin.H{
			"error": fmt.Errorf("cannot list devices: %v", err.Error()),
		})
	}

	c.JSON(http.StatusOK, response{
		RequestID: requestID,
		Payload: payload{
			UserID:  userID,
			Devices: devices,
		},
	})
}

func (h *yandexHandler) handleDevicesQuery(c *gin.Context) {
	requestID := getValueFromContext(c, "requestID")

	var devicesQuery models.DevicesQuery
	if err := c.ShouldBindQuery(&devicesQuery); err != nil {
		c.AbortWithStatusJSON(http.StatusBadRequest, gin.H{"error": "invalid request body"})
	}

	var deviceIDs []string
	for _, d := range devicesQuery.Devices {
		deviceIDs = append(deviceIDs, d.ID)
	}

	devices, err := h.devicesService.GetDevices(deviceIDs)
	if err != nil {
		log.Printf("cannot list devices: %s", err.Error())
		c.AbortWithStatusJSON(http.StatusInternalServerError, gin.H{
			"error": fmt.Errorf("cannot list devices: %v", err.Error()),
		})
	}

	c.JSON(http.StatusOK, response{
		RequestID: requestID,
		Payload: payload{
			Devices: devices,
		},
	})
}

func (h *yandexHandler) handleDevicesAction(c *gin.Context) {

}

func newYandexHandler(devicesService services.DevicesService) HandlerInterface {
	return &yandexHandler{
		devicesService: devicesService,
	}
}
