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
