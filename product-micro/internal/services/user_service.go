package services

import (
	"github.com/maxskaink/proyecto-microservicios/product-micro/internal/db/repositories"
	"github.com/maxskaink/proyecto-microservicios/product-micro/internal/dto"
)

// userService es la implementación de UserService.
type userService struct {
	userRepo repositories.IUserRepository
}

// NewUserService crea una nueva instancia de UserService.
func NewUserService(userRepo repositories.IUserRepository) IUserService {
	return &userService{userRepo: userRepo}
}

// GetUserByID obtiene un usuario por su ID.
func (s *userService) GetUserByID(id string) (dto.UserResponse, error) {
	// Llamar al repositorio para obtener el usuario
	user, err := s.userRepo.FindById(id)
	if err != nil {
		return dto.UserResponse{}, err
	}

	// Retornar el usuario como respuesta
	return *user, nil
}

func (s *userService) GetUserByUUID(id string) (dto.UserResponse, error) {
	user, err := s.userRepo.FindByUID(id)
	if err != nil {
		return dto.UserResponse{}, err
	}
	return *user, nil
}
