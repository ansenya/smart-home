package main

import (
	"auth-server/internal/app"
	"auth-server/internal/config"
	"github.com/gin-gonic/gin"
	"log"
)

func main() {
	cfg, err := config.Load()
	if err != nil {
		log.Fatalf("config load error: %v", err)
	}

	container, err := app.NewContainer(cfg)
	if err != nil {
		log.Fatalf("container build error: %v", err)
	}

	engine := gin.Default()

	container.Router.RegisterRoutes(engine)

	log.Printf("Starting server on %s...", cfg.Server.Port)
	if err := engine.Run(cfg.Server.Port); err != nil {
		log.Fatal(err)
	}
}
