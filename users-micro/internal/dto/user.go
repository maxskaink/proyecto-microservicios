package dto

import "github.com/maxskaink/proyecto-microservicios/users-micro/internal/domain"

type UserRequest struct {
	FirebaseUID string          `json:"firabseUID"`
	Email       string          `json:"email" binding:"required,email"`
	Name        string          `json:"name" binding:"required"`
	Rol         domain.UserRole `json:"rol"`
	Profile     *ProfileRequest `json:"profile"`
}

type UserResponse struct {
	ID          string          `json:"id"`
	FirebaseUID string          `json:"firabseUID"`
	Email       string          `json:"email"`
	Name        string          `json:"name"`
	Rol         domain.UserRole `json:"rol"`
	Profile     ProfileResponse `json:"profile"`
}
