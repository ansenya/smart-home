package models

type User struct {
	ID        string `gorm:"type:uuid;default:gen_random_uuid();primaryKey" json:"id"`
	Email     string `gorm:"uniqueIndex;not null" json:"email"`
	Password  string `gorm:"column:password;not null" json:"-"`
	Name      string `gorm:"column:name" json:"name,omitempty"`
	Confirmed bool   `gorm:"column:confirmed" json:"-"`
}
