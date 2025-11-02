package services

import (
	"context"

	"github.com/maxskaink/proyecto-microservicios/users-micro/internal/db/repositories"
	"github.com/maxskaink/proyecto-microservicios/users-micro/internal/domain"
	"github.com/maxskaink/proyecto-microservicios/users-micro/internal/dto"
	"github.com/maxskaink/proyecto-microservicios/users-micro/internal/messaging"
	"github.com/maxskaink/proyecto-microservicios/users-micro/pkg/logger"
)

// userService es la implementación de UserService.
type userService struct {
	userRepo  repositories.UserRepository
	publisher messaging.Publisher
}

// NewUserService crea una nueva instancia de UserService.
func NewUserService(userRepo repositories.UserRepository, publisher messaging.Publisher) UserService {
	return &userService{
		userRepo:  userRepo,
		publisher: publisher,
	}
}

func NewUserServiceWithoutPublisher(userRepo repositories.UserRepository) UserService {
	return &userService{
		userRepo: userRepo,
	}
}

// CreateUser crea un nuevo usuario.
func (s *userService) CreateUser(u dto.UserRequest, tenantID string) (dto.UserResponse, error) {
	// Validar tenant
	if tenantID == "" {
		return dto.UserResponse{}, domain.BadRequestError{Message: "Tenant requerido"}
	}

	// Validar los datos del usuario
	if u.Email == "" || u.Name == "" {
		return dto.UserResponse{}, domain.BadRequestError{Message: "El email y nombre no pueden estar vacios"}
	}

	if u.FirebaseUID == "" {
		return dto.UserResponse{}, domain.BadRequestError{Message: "El id de firebase no puede estar vacio"}
	}

	u.Rol = domain.UserRoleClient //Rol by default

	// Llamar al repositorio para crear el usuario
	userCreated, err := s.userRepo.Create(&u, tenantID)
	if err != nil {
		return dto.UserResponse{}, err
	}

	// Publicar evento de usuario creado
	if s.publisher != nil {
		if err := s.publisher.PublishUserCreated(context.Background(), *userCreated, tenantID); err != nil {
			logger.Error("Error al publicar evento de usuario creado: " + err.Error())
			// No retornamos error aquí para no afectar la operación principal
		} else {
			logger.Info("Evento de usuario creado publicado correctamente")
		}
	}

	// Retornar el usuario creado como respuesta
	return *userCreated, nil
}

// GetUserByID obtiene un usuario por su ID.
func (s *userService) GetUserByID(id string, tenantID string) (dto.UserResponse, error) {
	// Validar tenant
	if tenantID == "" {
		return dto.UserResponse{}, domain.BadRequestError{Message: "Tenant requerido"}
	}

	// Llamar al repositorio para obtener el usuario
	user, err := s.userRepo.FindByID(id, tenantID)
	if err != nil {
		return dto.UserResponse{}, err
	}

	// Retornar el usuario como respuesta
	return *user, nil
}

func (s *userService) GetUserByUUID(id string, tenantID string) (dto.UserResponse, error) {
	// Validar tenant
	if tenantID == "" {
		return dto.UserResponse{}, domain.BadRequestError{Message: "Tenant requerido"}
	}

	user, err := s.userRepo.FindByUUID(id, tenantID)
	if err != nil {
		return dto.UserResponse{}, err
	}
	return *user, nil
}

// UpdateUser actualiza un usuario existente. uid to know if has permission
func (s *userService) UpdateUser(id string, u dto.UserRequest, uid string, tenantID string) (dto.UserResponse, error) {
	// Validar tenant
	if tenantID == "" {
		return dto.UserResponse{}, domain.BadRequestError{Message: "Tenant requerido"}
	}

	user_to_update, err := s.userRepo.FindByID(id, tenantID)

	if err != nil {
		return dto.UserResponse{}, domain.NotFoundError{Message: "El usuario con id " + id + " no existe"}
	}

	if user_to_update.FirebaseUID != uid {
		return dto.UserResponse{}, domain.NotFoundError{Message: "El usuario no tiene permisos para editar el usuario"}
	}

	response, err := s.userRepo.Update(id, &u, tenantID)

	//Notificar por evento
	if err == nil && s.publisher != nil {
		if err := s.publisher.PublishUserUpdated(context.Background(), *response, tenantID); err != nil {
			logger.Error("Error al publicar evento de usuario actualizado: " + err.Error())
		} else {
			logger.Info("Evento de usuario actualizado publicado correctamente")
		}
	}

	return *response, err
}

// DeleteUser elimina un usuario por su ID.
func (s *userService) DeleteUser(id string, tenantID string) error {
	// Validar tenant
	if tenantID == "" {
		return domain.BadRequestError{Message: "Tenant requerido"}
	}

	//Is missing the validation of authorization
	return s.userRepo.Delete(id, tenantID)
}

// Update the rol of a user - uid of how are tring to change the rol
func (s *userService) UpdateRol(id string, rol domain.UserRole, uid_requester string, tenantID string) (dto.UserResponse, error) {
	// Validar tenant
	if tenantID == "" {
		return dto.UserResponse{}, domain.BadRequestError{Message: "Tenant requerido"}
	}

	//Validate the rol
	if !domain.IsValidUserRole(string(rol)) {
		return dto.UserResponse{}, domain.InvalidInputError{Message: "El rol debe ser: " + domain.StringValidRoles()}
	}
	//Validate athorization
	// Validar el rol del usuario solicitante (solo administradores pueden cambiar roles)
	requester, err := s.GetUserByUUID(uid_requester, tenantID)
	if err != nil {
		return dto.UserResponse{}, err
	}

	// Verificar si el usuario es administrador
	if string(requester.Rol) != string(domain.UserRoleAdmin) {
		return dto.UserResponse{}, domain.UnauthorizedError{Message: "Solo administradores pueden cambiar roles de usuario"}
	}
	//Validate if the user exist
	updated_user, err := s.userRepo.UpdateRol(id, string(rol), tenantID)
	if err != nil {
		return dto.UserResponse{}, err
	}

	if s.publisher != nil {
		if err := s.publisher.PublishUserUpdated(context.Background(), *updated_user, tenantID); err != nil {
			logger.Error("Error al publicar evento de usuario actualizado: " + err.Error())
		} else {
			logger.Info("Evento de usuario actualizado publicado correctamente")
		}
	}

	return *updated_user, err

}
