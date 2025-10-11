package models

import (
	"gorm.io/gorm"
	"time"
)

type OauthClient struct {
	ID          string `gorm:"primarykey"`
	Name        string
	Enabled     bool
	RedirectURI string
	UpdatedAt   time.Time
	CreatedAt   time.Time
	DeletedAt   gorm.DeletedAt `gorm:"index"`
}

func (OauthClient) TableName() string {
	return "oauth_clients"
}
