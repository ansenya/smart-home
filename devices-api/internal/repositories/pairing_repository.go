package repositories

import (
	"devices-api/internal/models"

	"github.com/google/uuid"
	"gorm.io/gorm"
)

type PairingRepository interface {
	FindManufacturedByMAC(mac string) (*models.ManufacturedDevice, error)
	RegisterDevice(userID uuid.UUID, request *models.ConfirmPairingRequest) error
}

type pairingRepository struct {
	db *gorm.DB
}

func NewPairingRepository(db *gorm.DB) PairingRepository {
	return &pairingRepository{
		db: db,
	}
}

func (p *pairingRepository) FindManufacturedByMAC(mac string) (*models.ManufacturedDevice, error) {
	var device models.ManufacturedDevice
	return &device, p.db.First(&device, "mac_address = ?", mac).Error
}

func (r *pairingRepository) RegisterDevice(userID uuid.UUID, request *models.ConfirmPairingRequest) error {
	return r.db.Transaction(func(tx *gorm.DB) error {
		return tx.Exec(`
			INSERT INTO devices (user_id, device_uid, mac_address, name, description, type)
			VALUES (?, ?, ?)
		`, userID, request.DeviceUID, request.MacAddress, request.Name, request.Description, request.Type).Error
	})
}
