package repositories

import (
	"devices-api/internal/models"

	"github.com/google/uuid"
	"gorm.io/gorm"
)

type PairingRepository interface {
	FindManufacturedByMAC(mac string) (*models.ManufacturedDevice, error)
	RegisterDevice(userID uuid.UUID, request *models.ConfirmPairingRequest) (uuid.UUID, error)
	DisablePreviouslyRegisteredDevice(uid string) error
}

type pairingRepository struct {
	db *gorm.DB
}

func NewPairingRepository(db *gorm.DB) PairingRepository {
	return &pairingRepository{
		db: db,
	}
}

func (r *pairingRepository) FindManufacturedByMAC(mac string) (*models.ManufacturedDevice, error) {
	var device models.ManufacturedDevice
	return &device, r.db.First(&device, "mac_address = ?", mac).Error
}

func (r *pairingRepository) RegisterDevice(userID uuid.UUID, request *models.ConfirmPairingRequest) (uuid.UUID, error) {
	device := &models.Device{
		UserID:      userID,
		DeviceUID:   request.DeviceUID,
		MacAddress:  request.MacAddress,
		Name:        request.Name,
		Description: request.Description,
		Type:        request.Type,
	}

	if err := r.db.Create(device).Error; err != nil {
		return uuid.Nil, err
	}

	return device.ID, nil
}

func (r *pairingRepository) DisablePreviouslyRegisteredDevice(uid string) error {
	return r.db.Delete(&models.Device{}).Where("device_uid = ?", uid).Error
}
