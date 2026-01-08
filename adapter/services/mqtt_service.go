package services

import (
	"adapter/models"
	"encoding/json"
	"fmt"
	"strings"

	mqtt "github.com/eclipse/paho.mqtt.golang"
)

type Component string

type mqttService struct {
	mqttClient mqtt.Client
}

// GetTopicName <user-id>/<device-id≥/<component>/<capability>/set
func (s *mqttService) GetTopicName(userID string, device *models.Device, componentName string) string {
	return strings.Join([]string{
		userID,
		device.ID,
		"capabilities",
		componentName,
		"set",
	}, "/")
}

func (s *mqttService) Publish(message any, topic string) error {
	payload, err := json.Marshal(message)
	if err != nil {
		return fmt.Errorf("failed to marshal message: %w", err)
	}

	token := s.mqttClient.Publish(topic, 0, true, payload)
	token.Wait()

	return token.Error()
}

func NewMqttService(mqttClient mqtt.Client) MqttService {
	return &mqttService{
		mqttClient: mqttClient,
	}
}
