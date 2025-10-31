package handlers

import (
	mqtt "github.com/eclipse/paho.mqtt.golang"
	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
)

type Router struct {
	healthHandler   HandlerInterface
	yandexV1Handler HandlerInterface
}

func NewRouter(db *gorm.DB, mqttClient *mqtt.Client) *Router {
	return &Router{
		healthHandler:   newHealthHandler(),
		yandexV1Handler: newYandexHandler(db, mqttClient),
	}
}

func (r *Router) RegisterRoutes(engine *gin.Engine) {
	r.healthHandler.RegisterRoutes(engine.Group(""))

	yandexV1Group := engine.Group("/yandex/v1.0")
	r.yandexV1Handler.RegisterRoutes(yandexV1Group)
}
