package models

import (
	"time"

	"github.com/google/uuid"
)

type ManufacturedDevice struct {
	ID         uuid.UUID `gorm:"primary_key"`
	Secret     string
	MacAddress string
	Registered bool
	CreatedAt  time.Time
}

type Device struct {
	ID          uuid.UUID `gorm:"primary_key;default:uuid_generate_v4()"`
	DeviceUID   string
	MacAddress  string
	UserID      uuid.UUID
	Name        string
	Description string
	Type        string
}
