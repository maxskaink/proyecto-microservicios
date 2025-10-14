package services

import (
	"fmt"

	"github.com/maxskaink/proyecto-microservicios/users-micro/internal/db/repositories"
	"github.com/maxskaink/proyecto-microservicios/users-micro/internal/domain"
	"github.com/maxskaink/proyecto-microservicios/users-micro/internal/dto"
)

// userService es la implementación de UserService.
type userService struct {
	userRepo repositories.UserRepository
}

// NewUserService crea una nueva instancia de UserService.
func NewUserService(userRepo repositories.UserRepository) UserService {
	return &userService{userRepo: userRepo}
}

// CreateUser crea un nuevo usuario.
func (s *userService) CreateUser(u dto.UserRequest) (dto.UserResponse, error) {
	// Validar los datos del usuario
	if u.Email == "" || u.Name == "" {
		return dto.UserResponse{}, fmt.Errorf("el email y el nombre son obligatorios")
	}

	u.Rol = domain.UserRoleClient
	// Llamar al repositorio para crear el usuario
	userCreated, err := s.userRepo.Create(&u)
	if err != nil {
		return dto.UserResponse{}, err
	}

	// Retornar el usuario creado como respuesta
	return dto.UserResponse{
		ID:    userCreated.ID, // Esto debería venir del repositorio
		Email: u.Email,
		Name:  u.Name,
	}, nil
}

// GetUserByID obtiene un usuario por su ID.
func (s *userService) GetUserByID(id string) (dto.UserResponse, error) {
	// Llamar al repositorio para obtener el usuario
	user, err := s.userRepo.FindByID(id)
	if err != nil {
		return dto.UserResponse{}, err
	}

	// Retornar el usuario como respuesta
	return *user, nil
}

func (s *userService) GetUserByUUID(id string) (dto.UserResponse, error) {
	user, err := s.userRepo.FindByUUID(id)

	if err != nil {
		return dto.UserResponse{}, err
	}
	return *user, nil
}

// UpdateUser actualiza un usuario existente.
func (s *userService) UpdateUser(id string, u dto.UserRequest) (dto.UserResponse, error) {
	// TODO: Implementar la lógica de negocio para actualizar un usuario.
	return dto.UserResponse{}, nil
}

// DeleteUser elimina un usuario por su ID.
func (s *userService) DeleteUser(id string) error {
	// TODO: Implementar la lógica de negocio para eliminar un usuario.
	return nil
}
