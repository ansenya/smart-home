package repository

import (
	"devices-api/models"
	"errors"
	"fmt"
	"gorm.io/gorm"
)

type capabilityRepo struct {
	db *gorm.DB
}

func (r *capabilityRepo) GetByDevice(deviceID string) ([]models.Capability, error) {
	var capabilities []models.Capability
	if err := r.db.Where("device_id = ?", deviceID).Find(&capabilities).Error; err != nil {
		return nil, err
	}
	return capabilities, nil
}

func (r *capabilityRepo) GetByID(id string) (*models.Capability, error) {
	var capability models.Capability
	if err := r.db.First(&capability, "id = ?", id).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, nil
		}
		return nil, fmt.Errorf("failed to get capability by id: %w", err)
	}
	return &capability, nil
}

func (r *capabilityRepo) UpdateState(capType string, state any) error {
	if err := r.db.Model(&models.Capability{}).
		Where("type = ?", capType).
		Update("state", state).Error; err != nil {
		return fmt.Errorf("failed to update capability state: %w", err)
	}
	return nil
}

func NewCapabilityRepo(db *gorm.DB) CapabilitiesRepository {
	return &capabilityRepo{
		db: db,
	}
}
