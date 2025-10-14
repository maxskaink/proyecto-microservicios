package domain

import "time"

type UserRole string

const (
	UserRoleAdmin    UserRole = "admin"
	UserRoleProducer UserRole = "producer"
	UserRoleClient   UserRole = "client"
)

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

func IsValidUserRole(role string) bool {
	switch UserRole(role) {
	case UserRoleAdmin, UserRoleProducer, UserRoleClient:
		return true
	default:
		return false
	}
}
