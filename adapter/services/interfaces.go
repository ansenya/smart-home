package services

import (
	"adapter/models"
	"encoding/json"
)

type StatusListenerService interface {
	StartListener()
}

type DevicesService interface {
	GetDevice(id string) (*models.Device, error)
	GetDevices(ids []string) ([]models.Device, error)
	GetUserDevices(userID string) ([]models.Device, error)
	CreateDevice(device *models.Device) error
	UpsertDevice(device *models.Device) error
	DeleteDevice(id string) error

	UpdateLastSeen(deviceID string) error
	UpdateCurrentState(deviceID, capability string, payload json.RawMessage) error
	UpdateProperty(deviceID, property string, payload json.RawMessage) error
}

type MqttService interface {
	GetTopicName(userID string, device *models.Device, componentName string) string
	Publish(message any, topic string) error
}
