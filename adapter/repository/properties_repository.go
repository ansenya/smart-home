package repository

import (
	"adapter/models"
	"log"

	"gorm.io/gorm"
)

type propertiesRepository struct {
	db *gorm.DB
}

func (p propertiesRepository) UpdateState(deviceID, property string, state *models.PropertyState) error {
	log.Println(property, state.Instance)
	return p.db.Model(&models.Property{}).
		Where("device_id = ?", deviceID).
		Where("parameters->>'instance' = ?", state.Instance).
		Update("state", state).Error
}

func (p propertiesRepository) ReplaceByDevice(deviceID string, props []models.Property) error {
	tx := p.db.Begin()

	if err := tx.Where("device_id = ?", deviceID).Delete(&models.Property{}).Error; err != nil {
		tx.Rollback()
		return err
	}

	for _, p := range props {
		p.DeviceID = deviceID
		if err := tx.Create(&p).Error; err != nil {
			tx.Rollback()
			return err
		}
	}

	return tx.Commit().Error
}

func NewPropertiesRepository(db *gorm.DB) PropertiesRepository {
	return &propertiesRepository{db: db}
}
