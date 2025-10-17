package services

import (
	"github.com/maxskaink/proyecto-microservicios/users-micro/internal/domain"
	"github.com/maxskaink/proyecto-microservicios/users-micro/internal/dto"
)

// UserService define la lógica de negocio para usuarios.
type UserService interface {
	CreateUser(u dto.UserRequest) (dto.UserResponse, error)
	GetUserByID(id string) (dto.UserResponse, error)
	GetUserByUUID(id string) (dto.UserResponse, error)
	UpdateUser(id string, u dto.UserRequest, uid string) (dto.UserResponse, error)
	UpdateRol(id string, rol domain.UserRole) (dto.UserResponse, error)
	DeleteUser(id string) error
}
