package dto

import (
	"time"
)

// UserDTO representa un usuario
type UserDTO struct {
	ID          string    `json:"id"`
	FirebaseUID string    `json:"firebaseUID"`
	Email       string    `json:"email"`
	Name        string    `json:"name"`
	Rol         string    `json:"rol"`
	CreatedAt   time.Time `json:"created_at"`
	UpdatedAt   time.Time `json:"updated_at"`
}
