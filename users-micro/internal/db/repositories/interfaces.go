package repositories

import (
	"github.com/maxskaink/proyecto-microservicios/users-micro/internal/dto"
)

// UserRepository define operaciones de persistencia para usuarios.
type UserRepository interface {
	Create(u *dto.UserRequest, tenantID string) (*dto.UserResponse, error)
	FindByID(id string, tenantID string) (*dto.UserResponse, error)
	FindByUUID(uuid string, tenantID string) (*dto.UserResponse, error)
	Update(id string, u *dto.UserRequest, tenantID string) (*dto.UserResponse, error)
	UpdateRol(id string, rol string, tenantID string) (*dto.UserResponse, error)
	Delete(id string, tenantID string) error
}
