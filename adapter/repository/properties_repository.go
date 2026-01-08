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

func NewPropertiesRepository(db *gorm.DB) PropertiesRepository {
	return &propertiesRepository{db: db}
}
