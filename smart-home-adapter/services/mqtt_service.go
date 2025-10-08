package services

import (
	"devices-api/models"
	"encoding/json"
	"fmt"
	mqtt "github.com/eclipse/paho.mqtt.golang"
	"strings"
)

type Component string

const (
	CapabilityComponent Component = "capabilities"
	PropertyComponent   Component = "properties"
)

// <user-id>/<device-id≥/<component>/<capability>/<action>
// <user-id>/<device-id>/<component>/<property>/<action>

type mqttService struct {
	mqttClient mqtt.Client
}

func (s *mqttService) GetTopicName(device *models.Device, component Component, componentName string, action string) string {
	return strings.Join([]string{
		device.UserID,
		device.ID,
		string(component),
		componentName,
		action,
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
