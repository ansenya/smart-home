package models

import "encoding/json"

type DevicesQuery struct {
	Devices []DeviceQuery `json:"devices"`
}

type DeviceQuery struct {
	ID         string          `json:"id"`
	CustomData json.RawMessage `json:"custom_data"`
}
