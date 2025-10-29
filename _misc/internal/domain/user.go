package domain

import "time"

type User struct {
	ID        string
	Email     string
	Password  string
	CreatedAt time.Time
}
