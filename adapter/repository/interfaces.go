package repository

import (
	"adapter/models"
	"encoding/json"
)

type DevicesRepository interface {
	GetByID(id string) (*models.Device, error)
	GetByIDs(ids []string) ([]models.Device, error)
	GetByUserID(userID string) ([]models.Device, error)
	Save(device *models.Device) error
	Update(device *models.Device) error
	UpdateLastSeen(id string) error
	Upsert(device *models.Device) error
	Delete(id string) error
}

type CapabilitiesRepository interface {
	GetByDevice(deviceID string) ([]models.Capability, error)
	GetByID(id string) (*models.Capability, error)
	UpdateState(deviceID, capability string, state json.RawMessage) error
}

type PropertiesRepository interface {
	UpdateState(deviceID, property string, state *models.PropertyState) error
}
