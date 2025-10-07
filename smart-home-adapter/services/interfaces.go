package services

import "devices-api/models"

type DevicesService interface {
	GetDevice(id string) (*models.Device, error)
	GetDevices(ids []string) ([]models.Device, error)
	GetUserDevices(userID string) ([]models.Device, error)
	CreateDevice(device *models.Device) error
	UpdateDevice(device *models.Device) error
	UpsertDevice(device *models.Device) error
	DeleteDevice(id string) error
}
