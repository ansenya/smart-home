package handlers

import (
	"device-service/internal/config"
	"log/slog"

	"github.com/gin-contrib/cors"
	"github.com/gin-gonic/gin"
)

type Router struct {
	engine *gin.Engine
	config *config.Container

	healthHandler  *healthHandler
	pairingHandler *pairingHandler

	log *slog.Logger
}

func NewRouter(container *config.Container) *Router {
	engine := gin.Default()
	router := Router{
		engine:         engine,
		config:         container,
		healthHandler:  newHealthHandler(),
		pairingHandler: newPairingHandler(),
		log:            container.Log,
	}
	router.configureCors()
	router.registerRoutes()
	return &router
}

func (r *Router) registerRoutes() {
	r.healthHandler.RegisterRoutes(r.engine)
	r.pairingHandler.RegisterRoutes(r.engine.Group("/devices/pair"))
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
		AllowMethods:     []string{"GET", "POST", "OPTIONS"},
		AllowHeaders:     []string{"Origin", "Content-Type", "Accept", "Authorization"},
		ExposeHeaders:    []string{"Content-Length"},
		AllowCredentials: true,
		MaxAge:           86400,
	}))
}
