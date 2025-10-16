package repositories

import (
	"github.com/maxskaink/proyecto-microservicios/users-micro/internal/dto"
)

// UserRepository define operaciones de persistencia para usuarios.
type UserRepository interface {
	Create(u *dto.UserRequest) (*dto.UserResponse, error)
	FindByID(id string) (*dto.UserResponse, error)
	FindByUUID(uuid string) (*dto.UserResponse, error)
	Update(id string, u *dto.UserRequest) (*dto.UserResponse, error)
	Delete(id string) error
}

// ProfileRepository define operaciones de persistencia para perfiles.
type ProfileRepository interface {
	Create(p *dto.ProfileRequest) error
	FindByUserID(userID string) (*dto.ProfileResponse, error)
	Update(p *dto.ProfileRequest) error
}
