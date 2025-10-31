package services

import (
	"devices-api/models"
	"devices-api/repository"
	"encoding/json"
	"gorm.io/gorm"
)

type devicesService struct {
	devicesRepository      repository.DevicesRepository
	capabilitiesRepository repository.CapabilitiesRepository
	propertiesRepository   repository.PropertiesRepository
}

func (r devicesService) GetDevice(id string) (*models.Device, error) {
	return r.devicesRepository.GetByID(id)
}

func (r devicesService) GetDevices(ids []string) ([]models.Device, error) {
	return r.devicesRepository.GetByIDs(ids)
}

func (r devicesService) GetUserDevices(userID string) ([]models.Device, error) {
	return r.devicesRepository.GetByUserID(userID)
}

func (r devicesService) CreateDevice(device *models.Device) error {
	return r.devicesRepository.Save(device)
}

func (r devicesService) UpsertDevice(device *models.Device) error {
	return r.devicesRepository.Upsert(device)
}

func (r devicesService) DeleteDevice(id string) error {
	return r.devicesRepository.Delete(id)
}

func (r devicesService) UpdateLastSeen(deviceID string) error {
	return r.devicesRepository.UpdateLastSeen(deviceID)
}

func (r devicesService) UpdateCurrentState(deviceID, capability string, payload json.RawMessage) error {
	return r.capabilitiesRepository.UpdateState(deviceID, capability, payload)
}

func (r devicesService) UpdateProperty(deviceID, property string, payload json.RawMessage) error {
	return r.capabilitiesRepository.UpdateState(deviceID, property, payload)
}

func NewDevicesService(db *gorm.DB) DevicesService {
	return &devicesService{
		devicesRepository:      repository.NewDevicesRepo(db),
		capabilitiesRepository: repository.NewCapabilityRepo(db),
		propertiesRepository:   repository.NewPropertiesRepository(db),
	}
}
