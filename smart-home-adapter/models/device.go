package models

import (
	"encoding/json"
)

type Device struct {
	ID          string          `gorm:"type:uuid;primary_key" json:"id"`
	MacAddress  string          `json:"mac_address"`
	UserID      string          `gorm:"column:user_id" json:"-"`
	Description string          `json:"description"`
	Name        string          `json:"name"`
	Room        string          `json:"room,omitempty"`
	Type        string          `json:"type"`
	CustomData  json.RawMessage `gorm:"column:custom_data" json:"custom_data,omitempty"`
	StatusInfo  json.RawMessage `gorm:"column:status_info" json:"status_info"`
	DeviceInfo  DeviceInfo      `gorm:"column:device_info" json:"device_info,omitempty"`

	Capabilities []Capability `gorm:"foreignKey:DeviceID" json:"capabilities,omitempty"`
	Properties   []Property   `gorm:"foreignKey:DeviceID" json:"properties,omitempty"`
}

type DeviceInfo struct {
	ID           string `json:"-"`
	Manufacturer string `json:"manufacturer"`
	Model        string `json:"model"`
	HwVersion    string `json:"hw_version,omitempty"`
	SwVersion    string `json:"sw_version,omitempty"`
	DeviceID     string `gorm:"type:uuid;unique" json:"-"`
}

type Capability struct {
	ID          string          `gorm:"type:uuid;primary_key" json:"id"`
	DeviceID    string          `gorm:"column:device_id" json:"-"`
	Type        string          `json:"type"`
	Retrievable bool            `json:"retrievable"`
	Reportable  bool            `json:"reportable"`
	Parameters  json.RawMessage `json:"parameters"`
	State       json.RawMessage `json:"state"`
}

type Property struct {
	ID          string          `gorm:"type:uuid;primary_key" json:"id"`
	DeviceID    string          `gorm:"column:device_id" json:"-"`
	Type        string          `json:"type"`
	Retrievable bool            `json:"retrievable"`
	Reportable  bool            `json:"reportable"`
	Parameters  json.RawMessage `json:"parameters"`
	State       json.RawMessage `json:"state"`
}
