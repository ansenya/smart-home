package models

import "encoding/json"

type DevicesActionRequest struct {
	Payload ActionPayloadRequest `json:"payload"`
}

type DevicesActionResponse struct {
	RequestId string                `json:"request_id"`
	Payload   ActionPayloadResponse `json:"payload"`
}

type ActionPayloadRequest struct {
	Devices []DeviceActionRequest `json:"devices"`
}

type ActionPayloadResponse struct {
	Devices []DeviceActionResponse `json:"devices"`
}

type DeviceActionRequest struct {
	ID           string          `json:"id"`
	CustomData   json.RawMessage `json:"custom_data"`
	Capabilities []Capability    `json:"capabilities"`
}

type DeviceActionResponse struct {
	ID           string                 `json:"id"`
	CustomData   json.RawMessage        `json:"custom_data"`
	Capabilities []CapabilityWithAction `json:"capabilities"`
}

type CapabilityWithAction struct {
	Type  string `json:"type"`
	State State  `json:"state"`
}

type State struct {
	Instance     string       `json:"instance"`
	ActionResult ActionResult `json:"action_result"`
}

type ActionResult struct {
	Status       string `json:"status,omitempty"`
	ErrorCode    string `json:"error_code,omitempty"`
	ErrorMessage string `json:"error_message,omitempty"`
}
