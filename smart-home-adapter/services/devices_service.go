package services

import (
	"devices-api/models"
	"devices-api/repository"
)

type devicesService struct {
	devicesRepository repository.DevicesRepository
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

func (r devicesService) UpdateDevice(device *models.Device) error {
	return r.devicesRepository.Update(device)
}

func (r devicesService) UpsertDevice(device *models.Device) error {
	return r.devicesRepository.Upsert(device)
}

func (r devicesService) DeleteDevice(id string) error {
	return r.devicesRepository.Delete(id)
}

func NewDevicesService(devicesRepository repository.DevicesRepository) DevicesService {
	return &devicesService{devicesRepository: devicesRepository}
}
