package handlers

import (
	"devices-api/middleware"
	"devices-api/models"
	"devices-api/services"
	"fmt"
	mqtt "github.com/eclipse/paho.mqtt.golang"
	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
	"log"
	"net/http"
	"time"
)

type yandexHandler struct {
	devicesService services.DevicesService
	mqttService    services.MqttService
}

func (h *yandexHandler) RegisterRoutes(g *gin.RouterGroup) {
	g.GET("", h.handleAliveStatus)

	protected := g.Group("/")
	protected.Use(middleware.YandexMiddleware(), middleware.JWTMiddleware())
	protected.POST("/user/unlink", h.handleUnlink)
	protected.GET("/user/devices", h.handleDevices)
	protected.POST("/user/devices/query", h.handleDevicesQuery)
	protected.POST("/user/devices/action", h.handleDevicesAction)
}

func (h *yandexHandler) handleAliveStatus(c *gin.Context) {
	c.Status(http.StatusOK)
}

func (h *yandexHandler) handleUnlink(c *gin.Context) {
	c.Status(http.StatusOK)
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
	for _, device := range devices {
		log.Println(device.LastSeen.String())
	}

	resp := models.YandexResponse{
		RequestID: requestID,
		Payload: models.Payload{
			UserID:  userID,
			Devices: devices,
		},
	}

	c.JSON(http.StatusOK, resp)
}

func (h *yandexHandler) handleDevicesQuery(c *gin.Context) {
	requestID := getValueFromContext(c, "requestID")

	var req struct {
		Devices []models.Device `json:"devices"`
	}
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid request body"})
		return
	}

	var deviceIDs []string
	for _, d := range req.Devices {
		deviceIDs = append(deviceIDs, d.ID)
	}

	devices, err := h.devicesService.GetDevices(deviceIDs)
	if err != nil {
		log.Printf("cannot list devices: %s", err.Error())
		c.JSON(http.StatusInternalServerError, gin.H{
			"error": fmt.Errorf("cannot list devices: %v", err.Error()),
		})
		return
	}
	for i := range devices {
		device := &devices[i]
		if time.Since(device.LastSeen) > 90*time.Second {
			device.ErrorCode = "DEVICE_UNREACHABLE"
			device.ErrorMessage = "устройство недоступно"
		}
	}

	c.JSON(http.StatusOK, models.YandexResponse{
		RequestID: requestID,
		Payload: models.Payload{
			Devices: devices,
		},
	})
}

func (h *yandexHandler) handleDevicesAction(c *gin.Context) {
	requestID := getValueFromContext(c, "requestID")
	userID := getValueFromContext(c, "userID")

	var req struct {
		Payload models.Payload
	}
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"error": "invalid req body",
		})
		return
	}

	// по девайсам из запроса
	for i := range req.Payload.Devices {
		device := &req.Payload.Devices[i]

		// по capabilities из запроса
		for j := range device.Capabilities {
			deviceCapability := &device.Capabilities[j]
			state := &deviceCapability.State

			topic := h.mqttService.GetTopicName(userID, device, deviceCapability.Type)
			if err := h.mqttService.Publish(state, topic); err != nil {
				log.Printf("cannot publish payload: %s", err.Error())
				state.ActionResult = &models.ActionResult{
					Status:       "ERROR",
					ErrorCode:    "500",
					ErrorMessage: fmt.Sprintf("cannot update state: %s", err.Error()),
				}
			} else {
				state.ActionResult = &models.ActionResult{
					Status: "DONE",
				}
			}
			deviceCapability.State.Value = nil
		}
	}

	c.JSON(http.StatusOK, models.YandexResponse{
		RequestID: requestID,
		Payload: models.Payload{
			Devices: req.Payload.Devices,
		},
	})
}

func newYandexHandler(db *gorm.DB, mqttClient *mqtt.Client) HandlerInterface {
	return &yandexHandler{
		devicesService: services.NewDevicesService(db),
		mqttService:    services.NewMqttService(*mqttClient),
	}
}
