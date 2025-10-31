package repository

import (
	"encoding/json"
	"fmt"
	"gorm.io/gorm"
)

type propertiesRepository struct {
	db *gorm.DB
}

func (p propertiesRepository) UpdateState(deviceID, property string, state json.RawMessage) error {
	return fmt.Errorf("not implemented")
}

func NewPropertiesRepository(db *gorm.DB) PropertiesRepository {
	return &propertiesRepository{db: db}
}
