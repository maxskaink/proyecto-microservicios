package domain

import "time"

type User struct {
	ID          string
	FirebaseUID string
	Email       string
	Name        string
	CreatedAt   time.Time
	UpdatedAt   time.Time
	DeletedAt   *time.Time
}

type Profile struct {
	ID        string
	UserID    string
	Address   string
	Phone     string
	AvatarURL string
	CreatedAt time.Time
	UpdatedAt time.Time
}
