package models

import "github.com/google/uuid"

type ConfirmPairingRequest struct {
	Code      string    `json:"code" binding:"required"`
	DeviceUID uuid.UUID `json:"device_uid" binding:"required"`
}
