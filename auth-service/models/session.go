package models

import "time"

type Session struct {
	ID         string     `gorm:"type:uuid;default:gen_random_uuid();primaryKey" json:"id"`
	UserID     string     `json:"user_id"`
	LastActive time.Time  `json:"last_active"`
	CreatedAt  time.Time  `json:"created_at"`
	ExpiresAt  *time.Time `json:"expires_at"`
}

func (*Session) TableName() string {
	return "oauth_sessions"
}
