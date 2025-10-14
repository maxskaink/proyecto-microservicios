package services

import "github.com/maxskaink/proyecto-microservicios/users-micro/internal/dto"

// UserService define la lógica de negocio para usuarios.
type UserService interface {
	CreateUser(u dto.UserRequest) (dto.UserResponse, error)
	GetUserByID(id string) (dto.UserResponse, error)
	GetUserByUUID(id string) (dto.UserResponse, error)
	UpdateUser(id string, u dto.UserRequest) (dto.UserResponse, error)
	DeleteUser(id string) error
}

// ProfileService define la lógica de negocio para perfiles de usuario.
type ProfileService interface {
	CreateProfile(p dto.ProfileRequest) (dto.ProfileResponse, error)
	GetProfileByUserID(userID string) (dto.ProfileResponse, error)
	UpdateProfile(userID string, p dto.ProfileRequest) (dto.ProfileResponse, error)
}
