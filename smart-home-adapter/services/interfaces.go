package services

import (
	"devices-api/models"
	"encoding/json"
)

type DevicesService interface {
	GetDevice(id string) (*models.Device, error)
	GetDevices(ids []string) ([]models.Device, error)
	GetUserDevices(userID string) ([]models.Device, error)
	CreateDevice(device *models.Device) error
	UpdateDevice(device *models.Device) error
	UpsertDevice(device *models.Device) error
	DeleteDevice(id string) error

	UpdateCapabilityState(capID string, state json.RawMessage) error
	UpdateCapabilitiesState(capID []string, state []any) error
}

type MqttService interface {
	GetTopicName(device *models.Device, component Component, componentName string, action string) string
	Publish(message any, topic string) error
}
