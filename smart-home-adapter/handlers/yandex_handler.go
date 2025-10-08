package handlers

import (
	"devices-api/handlers/middleware"
	"devices-api/models"
	"devices-api/services"
	"encoding/json"
	"fmt"
	"github.com/gin-gonic/gin"
	"log"
	"net/http"
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
	mqttService    services.MqttService
}

func (h *yandexHandler) RegisterRoutes(g *gin.RouterGroup) {
	g.GET("", h.handleAliveStatus)

	protected := g.Group("/")
	protected.Use(middleware.YandexMiddleware()) // должен быть подтвержденный акк
	protected.POST("/user/unlink", h.handleUnlink)
	protected.GET("/user/devices", h.handleDevices)
	protected.POST("/user/devices/query", h.handleDevicesQuery)
	protected.POST("/user/devices/action", h.handleDevicesAction)
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

	var devicesQuery struct {
		Devices []models.Device `json:"devices"`
	}
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
	requestID := getValueFromContext(c, "requestID")

	var req struct {
		Payload struct {
			Devices []models.Device `json:"devices"`
		} `json:"payload"`
	}
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"error": "invalid req body",
		})
		return
	}

	resp := response{
		RequestID: requestID,
	}

	for i := range req.Payload.Devices {
		d := &req.Payload.Devices[i]
		for j := range d.Capabilities {
			cp := &d.Capabilities[j]

			var state models.State
			if err := json.Unmarshal(cp.State, &state); err != nil {
				log.Printf("cannot unmarshal state: %s", err.Error())
			}

			err := h.devicesService.UpdateCapabilityState(cp.ID, cp.State)
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

			m, _ := json.Marshal(state)
			cp.State = m
		}
		resp.Payload.Devices = append(resp.Payload.Devices, *d)
	}

	c.JSON(http.StatusOK, resp)
}

func newYandexHandler(devicesService services.DevicesService, mqttService services.MqttService) HandlerInterface {
	return &yandexHandler{
		devicesService: devicesService,
		mqttService:    mqttService,
	}
}
