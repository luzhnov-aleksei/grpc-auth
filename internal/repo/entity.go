package repo

import "time"

type User struct {
	Email       string
	Username    string
	HashPass    string
	FirstName   string
	LastName    string
	LastLoginAt time.Time
	CreatedAt   time.Time
	UpdatedAt   time.Time
}
