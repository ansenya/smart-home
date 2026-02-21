package handlers

import (
	"llm-service/internal/config"
	"llm-service/internal/middleware"
	"llm-service/internal/repositories"
	"llm-service/internal/services"
	"log"
	"log/slog"

	"github.com/gin-contrib/cors"
	"github.com/gin-gonic/gin"
)

type Handler interface {
	RegisterRoutes(rg *gin.RouterGroup)
}
type Router struct {
	engine *gin.Engine
	cfg    *config.Container
	log    *slog.Logger

	HealthHandler   Handler
	ChatHandler     Handler
	MessagesHandler Handler
}

func (r *Router) registerRoutes(repos *repositories.Container) {
	chatsGroup := r.engine.Group("/chats")
	chatsGroup.Use(middleware.SessionAuth(repos.SessionRepository))

	r.ChatHandler.RegisterRoutes(chatsGroup)

	messagesGroup := chatsGroup.Group("/:id")
	r.MessagesHandler.RegisterRoutes(messagesGroup)

	// Health
	r.HealthHandler.RegisterRoutes(r.engine.Group(""))
}

func (r *Router) configureCors() {
	log.Printf("initializing cors")

	r.engine.Use(cors.New(cors.Config{
		AllowOrigins: []string{
			"https://smarthome.hipahopa.ru",
		},
		AllowMethods:     []string{"GET", "POST", "OPTIONS"},
		AllowHeaders:     []string{"Origin", "Content-Type", "Accept", "Authorization"},
		ExposeHeaders:    []string{"Content-Length"},
		AllowCredentials: true,
		MaxAge:           86400,
	}))
}

func (r *Router) Run() error {
	r.log.Info("starting server")
	return r.engine.Run(r.cfg.Server.Port)
}

func NewRouter(cfg *config.Container, svcs *services.Container, repos *repositories.Container) *Router {
	router := &Router{
		engine: gin.New(),
		log:    cfg.Log,
		cfg:    cfg,

		HealthHandler:   NewHealthHandler(),
		ChatHandler:     NewChatHandler(svcs),
		MessagesHandler: NewMessagesHandler(svcs),
	}

	router.configureCors()
	router.registerRoutes(repos)

	return router
}
