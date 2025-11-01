package main

import (
	"log"
	"registry/config"
	"registry/handlers"
)

func main() {
	container := config.NewConfig()
	router := handlers.NewRouter(container)
	if err := router.Run(); err != nil {
		log.Fatalf("failed to start: %s", err)
	}
}
