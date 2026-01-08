package repositories

import (
	"devices-api/internal/models"

	"github.com/google/uuid"
	"gorm.io/gorm"
)

type PairingRepository interface {
	FindManufacturedByMAC(mac string) (*models.ManufacturedDevice, error)
	RegisterDevice(deviceID uuid.UUID, userID uuid.UUID) error
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

func (r *pairingRepository) RegisterDevice(deviceID, userID uuid.UUID) error {
	return r.db.Transaction(func(tx *gorm.DB) error {
		if err := tx.Exec(`
			INSERT INTO devices (id, user_id)
			VALUES (?, ?)
		`, deviceID, userID).Error; err != nil {
			return err
		}

		if err := tx.Exec(`
			UPDATE manufactured_devices
			SET registered = true
			WHERE id = ?
		`, deviceID).Error; err != nil {
			return err
		}

		return nil
	})
}
