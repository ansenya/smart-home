package handlers

import (
	"devices-api/middleware"
	"devices-api/models"
	"devices-api/services"
	"encoding/json"
	"fmt"
	"github.com/gin-gonic/gin"
	"log"
	"net/http"
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
		c.AbortWithStatusJSON(http.StatusBadRequest, gin.H{"error": "invalid request body"})
	}

	var deviceIDs []string
	for _, d := range req.Devices {
		deviceIDs = append(deviceIDs, d.ID)
	}

	devices, err := h.devicesService.GetDevices(deviceIDs)
	if err != nil {
		log.Printf("cannot list devices: %s", err.Error())
		c.AbortWithStatusJSON(http.StatusInternalServerError, gin.H{
			"error": fmt.Errorf("cannot list devices: %v", err.Error()),
		})
	}

	m, _ := json.Marshal(devices)
	log.Println(string(m))

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

	for i := range req.Payload.Devices {
		d := &req.Payload.Devices[i]
		for j := range d.Capabilities {
			cp := &d.Capabilities[j]
			state := &cp.State
			m, _ := json.Marshal(state)

			topic := h.mqttService.GetTopicName(userID, d, services.CapabilityComponent, cp.Type)
			if err := h.mqttService.Publish(state, topic); err != nil {
				log.Printf("cannot publish payload: %s", err.Error())
			}

			err := h.devicesService.UpdateCapabilityState(cp.Type, m)
			if err != nil {
				state.ActionResult = models.ActionResult{
					Status:       "ERROR",
					ErrorCode:    "500",
					ErrorMessage: fmt.Sprintf("cannot update state: %s", err.Error()),
				}
			} else {
				state.ActionResult = models.ActionResult{
					Status: "DONE",
				}
			}
			cp.State.Value = nil
		}
	}

	c.JSON(http.StatusOK, models.YandexResponse{
		RequestID: requestID,
		Payload: models.Payload{
			Devices: req.Payload.Devices,
		},
	})
}

func newYandexHandler(devicesService services.DevicesService, mqttService services.MqttService) HandlerInterface {
	return &yandexHandler{
		devicesService: devicesService,
		mqttService:    mqttService,
	}
}
