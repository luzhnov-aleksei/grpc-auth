package repo

import "time"

type User struct {
	ID          int
	Email       string
	Username    string
	PassHash    string
	FirstName   string
	LastName    string
	LastLoginAt time.Time
	CreatedAt   time.Time
	UpdatedAt   time.Time
}
