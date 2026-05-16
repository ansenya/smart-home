package smarthome

import (
	"encoding/json"

	"gorm.io/gorm"
)

type Device struct {
	ID           string         `gorm:"type:uuid;primaryKey"`
	UserID       string         `gorm:"column:user_id"`
	Name         string         `gorm:"column:name"`
	Description  string         `gorm:"column:description"`
	Room         string         `gorm:"column:room"`
	Type         string         `gorm:"column:type"`
	Capabilities []Capability   `gorm:"foreignKey:DeviceID"`
	Properties   []Property     `gorm:"foreignKey:DeviceID"`
	DeletedAt    gorm.DeletedAt `gorm:"index"`
}

func (Device) TableName() string { return "devices" }

type Capability struct {
	ID         string          `gorm:"type:uuid;primaryKey"`
	DeviceID   string          `gorm:"column:device_id"`
	Type       string          `gorm:"column:type"`
	Parameters json.RawMessage `gorm:"column:parameters"`
	State      json.RawMessage `gorm:"column:state"`
}

func (Capability) TableName() string { return "capabilities" }

type Property struct {
	ID         string          `gorm:"type:uuid;primaryKey"`
	DeviceID   string          `gorm:"column:device_id"`
	Type       string          `gorm:"column:type"`
	Parameters json.RawMessage `gorm:"column:parameters"`
	State      json.RawMessage `gorm:"column:state"`
}

func (Property) TableName() string { return "properties" }
