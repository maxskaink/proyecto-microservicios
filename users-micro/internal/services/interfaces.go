package services

import (
	"github.com/maxskaink/proyecto-microservicios/users-micro/internal/domain"
	"github.com/maxskaink/proyecto-microservicios/users-micro/internal/dto"
)

// UserService define la lógica de negocio para usuarios.
type UserService interface {
	CreateUser(u dto.UserRequest, tenantID string) (dto.UserResponse, error)
	GetUserByID(id string, tenantID string) (dto.UserResponse, error)
	GetUserByUUID(id string, tenantID string) (dto.UserResponse, error)
	GetProducerByID(id string, tenantID string) (dto.UserResponse, error)
	ListUsers(tenantID string, uuid string) ([]dto.UserResponse, error)
	UpdateUser(id string, u dto.UserRequest, uid string, tenantID string) (dto.UserResponse, error)
	UpdateRol(id string, rol domain.UserRole, uid string, tenantID string) (dto.UserResponse, error)
	DeleteUser(id string, tenantID string) error
}
