package handlers

import (
	"log/slog"
	"os"
	"panel-api/internal/config"
	"strings"

	"github.com/gin-contrib/cors"
	"github.com/gin-gonic/gin"
)

const SessionIDName = "sid"

var DomainName = func() string {
	if v := os.Getenv("PANEL_COOKIE_DOMAIN"); v != "" {
		return v
	}
	return ".smarthome.hipahopa.ru"
}()

type Router struct {
	engine *gin.Engine
	config *config.Container

	healthHandler *healthHandler
	usersHandler  *usersHandler

	log *slog.Logger
}

func NewRouter(container *config.Container) *Router {
	engine := gin.Default()
	router := Router{
		engine:        engine,
		config:        container,
		healthHandler: newHealthHandler(),
		usersHandler:  newUsersHandler(container),
		log:           container.Log,
	}
	router.configureCors()
	router.registerRoutes()
	return &router
}

func (r *Router) registerRoutes() {
	r.healthHandler.RegisterRoutes(r.engine)
	r.usersHandler.RegisterRoutes(r.engine.Group("/panel/v1/users"))
}

func (r *Router) Run() error {
	return r.engine.Run(r.config.Server.Port)
}

func (r *Router) configureCors() {
	r.log.Info("initializing cors")

	origins := []string{"https://smarthome.hipahopa.ru", "http://localhost:5173"}
	if extra := os.Getenv("CORS_ALLOWED_ORIGINS"); extra != "" {
		origins = append(origins, strings.Split(extra, ",")...)
	}

	r.engine.Use(cors.New(cors.Config{
		AllowOrigins:     origins,
		AllowMethods:     []string{"GET", "POST", "OPTIONS"},
		AllowHeaders:     []string{"Origin", "Content-Type", "Accept", "Authorization"},
		ExposeHeaders:    []string{"Content-Length"},
		AllowCredentials: true,
		MaxAge:           86400,
	}))
}
