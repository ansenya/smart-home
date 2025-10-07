package models

import (
	"gorm.io/gorm"
)

type User struct {
	ID        string `gorm:"type:uuid;primaryKey" json:"id"`
	Email     string `gorm:"uniqueIndex;not null" json:"email"`
	Password  string `gorm:"column:password;not null"`
	Name      string `gorm:"column:name"`
	Confirmed bool   `gorm:"column:confirmed"`
	gorm.Model
}

type RegistrationUser struct {
	Email    string `json:"email" binding:"required"`
	Password string `json:"password" binding:"required"`
	Name     string `json:"name"`
}

type LoginUser struct {
	Email    string `json:"email" binding:"required"`
	Password string `json:"password" binding:"required"`
}
