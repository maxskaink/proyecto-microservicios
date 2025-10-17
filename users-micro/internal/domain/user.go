package domain

import (
	"strings"
	"time"
)

type UserRole string

const (
	UserRoleAdmin    UserRole = "admin"
	UserRoleProducer UserRole = "producer"
	UserRoleClient   UserRole = "client"
)

var ValidRoles = []UserRole{
	UserRoleAdmin,
	UserRoleProducer,
	UserRoleClient,
}

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

func StringValidRoles() string {
	roles := make([]string, len(ValidRoles))
	for i, role := range ValidRoles {
		roles[i] = string(role)
	}
	return strings.Join(roles, ", ")
}
