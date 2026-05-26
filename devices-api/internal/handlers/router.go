package handlers

import (
	"devices-api/internal/config"
	"devices-api/internal/middleware"
	"devices-api/internal/repositories"
	"devices-api/internal/services"
	"log/slog"

	"github.com/gin-contrib/cors"
	"github.com/gin-gonic/gin"
)

type Router struct {
	engine *gin.Engine
	config *config.Container

	healthHandler  *healthHandler
	pairingHandler *pairingHandler
	devicesHandler *devicesHandler
	streamHandler  *streamHandler

	log *slog.Logger
}

func NewRouter(
	container *config.Container,
	repo repositories.SessionRepository,
	pairing services.PairingService,
	devices services.DevicesService,
	stream services.StreamService,
) *Router {
	engine := gin.Default()
	router := Router{
		engine:         engine,
		config:         container,
		healthHandler:  newHealthHandler(),
		pairingHandler: newPairingHandler(pairing),
		devicesHandler: newDevicesHandler(devices),
		streamHandler:  newStreamHandler(stream),
		log:            container.Log,
	}
	router.configureCors()
	router.registerRoutes(repo)
	return &router
}

func (r *Router) registerRoutes(repo repositories.SessionRepository) {
	r.healthHandler.RegisterRoutes(r.engine)

	pairingGroup := r.engine.Group("/devices/pairing")
	pairingGroup.Use(middleware.SessionAuth(repo))
	devicesPairingGroup := r.engine.Group("/devices/pairing")
	r.pairingHandler.RegisterRoutes(pairingGroup, devicesPairingGroup)

	devicesGroup := r.engine.Group("/devices")
	devicesGroup.Use(middleware.SessionAuth(repo))
	r.devicesHandler.RegisterRoutes(devicesGroup)
	r.streamHandler.RegisterRoutes(devicesGroup)
}

func (r *Router) Run() error {
	return r.engine.Run(r.config.Server.Port)
}

func (r *Router) configureCors() {
	r.log.Info("initializing cors")
	r.engine.Use(cors.New(cors.Config{
		AllowOrigins: []string{
			"https://smarthome.hipahopa.ru",
			"http://localhost:5173",
		},
		AllowMethods:     []string{"GET", "POST", "PUT", "DELETE", "OPTIONS"},
		AllowHeaders:     []string{"Origin", "Content-Type", "Accept", "Authorization"},
		ExposeHeaders:    []string{"Content-Length"},
		AllowCredentials: true,
		MaxAge:           86400,
	}))
}
