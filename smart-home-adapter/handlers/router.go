package handlers

import (
	"devices-api/services"
	"github.com/gin-gonic/gin"
)

type Router struct {
	healthHandler   HandlerInterface
	yandexV1Handler HandlerInterface
}

func NewRouter(devicesService services.DevicesService, mqttService services.MqttService) *Router {
	return &Router{
		healthHandler:   newHealthHandler(),
		yandexV1Handler: newYandexHandler(devicesService, mqttService),
	}
}

func (r *Router) RegisterRoutes(engine *gin.Engine) {
	yandexV1Group := engine.Group("/yandex/v1.0")
	r.yandexV1Handler.RegisterRoutes(yandexV1Group)

	r.healthHandler.RegisterRoutes(engine.Group(""))
}
