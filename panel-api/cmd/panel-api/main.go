package panel_api

import (
	"fmt"
	"panel-api/internal/config"
	"panel-api/internal/handlers"
)

func main() {
	// config
	container := config.NewConfig()
	container.Log.Info("notification service started")

	router := handlers.NewRouter(container)
	if err := router.Run(); err != nil {
		container.Log.Error(fmt.Sprintf("failed to start: %s", err))
	}
}
