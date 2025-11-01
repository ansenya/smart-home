package models

import (
	"encoding/json"
	"time"
)

// <user-id>/<device-id≥/capabilities/<capability>/set
// <user-id>/<device-id≥/capabilities/<capability>/state
// <user-id>/<device-id>/properties/<property>/state
// <user-id>/<device-id>/state

type YandexResponse struct {
	RequestID string  `json:"requestId"`
	Payload   Payload `json:"payload"`
}

type Payload struct {
	UserID  string   `json:"user_id,omitempty"`
	Devices []Device `json:"devices"`
}

type Device struct {
	ID           string          `gorm:"type:uuid;primary_key" json:"id"`
	MacAddress   string          `json:"-"`
	UserID       string          `gorm:"column:user_id" json:"-"`
	Name         string          `json:"name,omitempty"`
	Description  string          `json:"description,omitempty"`
	Room         string          `json:"room,omitempty"`
	Type         string          `json:"type,omitempty"`
	CustomData   json.RawMessage `gorm:"column:custom_data" json:"custom_data,omitempty"`
	StatusInfo   StatusInfo      `gorm:"type:jsonb;serializer:json;column:status_info" json:"status_info,omitempty"`
	DeviceInfo   DeviceInfo      `gorm:"type:jsonb;serializer:json;column:device_info" json:"device_info,omitempty"`
	Capabilities []Capability    `gorm:"foreignKey:DeviceID" json:"capabilities,omitempty"`
	Properties   []Property      `gorm:"foreignKey:DeviceID" json:"properties,omitempty"`

	LastSeen time.Time `gorm:"column:last_seen" json:"-"`

	ErrorCode    string `json:"error_code,omitempty"`
	ErrorMessage string `json:"error_message,omitempty"`
}

type StatusInfo struct {
	Reportable bool `json:"reportable"`
}

type Capability struct {
	ID          string          `gorm:"type:uuid;primary_key" json:"-"`
	DeviceID    string          `gorm:"column:device_id" json:"-"`
	Type        string          `json:"type"`
	Retrievable bool            `json:"retrievable"`
	Reportable  bool            `json:"reportable"`
	Parameters  json.RawMessage `json:"parameters,omitempty"`
	State       State           `gorm:"type:jsonb;serializer:json" json:"state,omitempty"`
}

type Property struct {
	ID          string `gorm:"type:uuid;primary_key" json:"-"`
	DeviceID    string `gorm:"column:device_id" json:"-"`
	Type        string
	Retrievable bool
	Reportable  bool
	Parameters  json.RawMessage
	State       State `gorm:"type:jsonb;serializer:json" json:"state,omitempty"`
}

type DeviceInfo struct {
	Manufacturer string `json:"manufacturer"`
	Model        string `json:"model"`
	HwVersion    string `json:"hw_version"`
	SwVersion    string `json:"sw_version"`
}

type State struct {
	Instance     string          `json:"instance"`
	Value        json.RawMessage `json:"value,omitempty"`
	ActionResult *ActionResult   `json:"action_result,omitempty"`
}

type ActionResult struct {
	Status       string `json:"status,omitempty"`
	ErrorCode    string `json:"error_code,omitempty"`
	ErrorMessage string `json:"error_message,omitempty"`
}

type PropertyState struct {
	Instance string `json:"instance"`
	Value    any    `json:"value"`
}
