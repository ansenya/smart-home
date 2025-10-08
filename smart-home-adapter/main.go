package main

import (
	"devices-api/handlers"
	"devices-api/repository"
	"devices-api/services"
	"devices-api/storage"
	"devices-api/utils"
	"fmt"
	mqtt "github.com/eclipse/paho.mqtt.golang"
	"log"
	"time"

	"github.com/gin-gonic/gin"
)

var (
	Port        = fmt.Sprintf(":%s", utils.GetEnv("PORT", "8080"))
	PostgresUrl = utils.GetEnv("POSTGRES_URL", "host=localhost user=user password=password dbname=smart-home port=5432 sslmode=disable")
	MqttUrl     = utils.GetEnv("MQTT_URL", "tcp://localhost:1883")
)

func main() {
	database, err := storage.ConnectPostgres(PostgresUrl)
	if err != nil {
		log.Fatalf("failed to connect to postgres: %v", err)
	}

	mqttClient, err := storage.ConnectMQTT(createMQTTOptions())
	if err != nil {
		log.Fatalf("failed to connect to mqtt: %v", err)
	}

	// repositories
	devicesRepository := repository.NewDevicesRepo(database)
	capabilitiesRepo := repository.NewCapabilityRepo(database)

	// services
	devicesService := services.NewDevicesService(devicesRepository, capabilitiesRepo)
	mqttService := services.NewMqttService(mqttClient)

	engine := gin.Default()

	router := handlers.NewRouter(devicesService, mqttService)
	router.RegisterRoutes(engine)

	if err := engine.Run(Port); err != nil {
		log.Fatalf("failed to start server: %v", err)
	}
}

func createMQTTOptions() *mqtt.ClientOptions {
	opts := mqtt.NewClientOptions()
	opts.AddBroker(MqttUrl)
	opts.SetUsername("") // Optional
	opts.SetPassword("") // Optional
	opts.SetKeepAlive(2 * time.Second)
	opts.SetPingTimeout(1 * time.Second)
	opts.OnConnect = func(c mqtt.Client) {
		fmt.Println("Connected to MQTT broker")
	}
	opts.OnConnectionLost = func(c mqtt.Client, err error) {
		log.Fatalf("Connection to MQTT broker lost: %v\n", err)
	}

	return opts
}
