package models

type User struct {
	ID        string `gorm:"type:uuid;default:gen_random_uuid();primaryKey" json:"-"`
	Email     string `gorm:"uniqueIndex;not null" json:"email"`
	Password  string `gorm:"column:password;not null" json:"-"`
	Name      string `gorm:"column:name" json:"name,omitempty"`
	Confirmed bool   `gorm:"column:confirmed" json:"-"`
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
