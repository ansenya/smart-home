package models

import (
	"encoding/json"
	"time"

	"github.com/google/uuid"
	"gorm.io/gorm"
)

type DeviceWithRelations struct {
	ID           uuid.UUID       `gorm:"type:uuid;primaryKey;default:uuid_generate_v4()" json:"id"`
	DeviceUID    string          `gorm:"column:device_uid" json:"device_uid"`
	MacAddress   string          `gorm:"column:mac_address" json:"-"`
	UserID       uuid.UUID       `gorm:"column:user_id" json:"-"`
	Name         string          `json:"name"`
	Description  string          `json:"description"`
	Room         string          `json:"room"`
	Type         string          `json:"type"`
	StatusInfo   json.RawMessage `gorm:"column:status_info" json:"status_info,omitempty"`
	CustomData   json.RawMessage `gorm:"column:custom_data" json:"custom_data,omitempty"`
	DeviceInfo   json.RawMessage `gorm:"column:device_info" json:"device_info,omitempty"`
	LastSeen     time.Time       `gorm:"column:last_seen" json:"last_seen"`
	CreatedAt    time.Time       `gorm:"column:created_at" json:"created_at"`
	UpdatedAt    time.Time       `gorm:"column:updated_at" json:"updated_at"`
	DeletedAt    gorm.DeletedAt  `gorm:"index" json:"-"`
	Capabilities []Capability    `gorm:"foreignKey:DeviceID" json:"capabilities"`
	Properties   []Property      `gorm:"foreignKey:DeviceID" json:"properties"`
}

func (DeviceWithRelations) TableName() string { return "devices" }

type Capability struct {
	ID          uuid.UUID       `gorm:"type:uuid;primaryKey;default:uuid_generate_v4()" json:"id"`
	DeviceID    uuid.UUID       `gorm:"column:device_id" json:"-"`
	Type        string          `json:"type"`
	Retrievable bool            `json:"retrievable"`
	Reportable  bool            `json:"reportable"`
	Parameters  json.RawMessage `json:"parameters,omitempty"`
	State       json.RawMessage `json:"state,omitempty"`
}

func (Capability) TableName() string { return "capabilities" }

type Property struct {
	ID          uuid.UUID       `gorm:"type:uuid;primaryKey;default:uuid_generate_v4()" json:"id"`
	DeviceID    uuid.UUID       `gorm:"column:device_id" json:"-"`
	Type        string          `json:"type"`
	Retrievable bool            `json:"retrievable"`
	Reportable  bool            `json:"reportable"`
	Parameters  json.RawMessage `json:"parameters,omitempty"`
	State       json.RawMessage `json:"state,omitempty"`
}

func (Property) TableName() string { return "properties" }

type UpdateDeviceRequest struct {
	Name        *string `json:"name"`
	Description *string `json:"description"`
	Room        *string `json:"room"`
}

type SetCapabilityRequest struct {
	Instance string          `json:"instance"`
	Value    json.RawMessage `json:"value" binding:"required"`
}

type DeviceListResponse struct {
	Devices []DeviceWithRelations `json:"devices"`
}
