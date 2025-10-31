package handlers

import (
	"encoding/json"
	"fmt"

	"github.com/maxskaink/proyecto-microservicios/product-micro/internal/db/repositories"
	"github.com/maxskaink/proyecto-microservicios/product-micro/internal/dto"
	"github.com/maxskaink/proyecto-microservicios/product-micro/internal/messaging/events"
	"github.com/maxskaink/proyecto-microservicios/product-micro/pkg/logger"
)

// UserEventHandler maneja eventos de usuario
type UserEventHandler struct {
	UserRepository repositories.IUserRepository
}

// HandleUserCreated procesa el evento de usuario creado
func (h *UserEventHandler) HandleUserCreated(data []byte) error {
	var user dto.UserRequest

	if err := json.Unmarshal(data, &user); err != nil {
		logger.Error(fmt.Sprintf("Error al deserializar evento de usuario creado: %v", err))
		return err
	}

	userCreated, err := h.UserRepository.CreateUser(user)

	if err != nil {
		logger.Error("No se ha podido crear el usuario con id " + user.ID)
	}

	logger.Info(fmt.Sprintf("Usuario creado: ID=%s, Email=%s, Nombre=%s, Rol=%s",
		userCreated.ID, userCreated.Email, userCreated.Name, userCreated.Rol))

	return nil
}

// HandleUserDeleted procesa el evento de usuario eliminado
func (h *UserEventHandler) HandleUserDeleted(data []byte) error {
	var event events.UserDeletedEvent

	if err := json.Unmarshal(data, &event); err != nil {
		logger.Error(fmt.Sprintf("Error al deserializar evento de usuario eliminado: %v", err))
		return err
	}

	err := h.UserRepository.DeleteUser(event.ID)
	if err != nil {
		logger.Error("No se ha podido eliminar el usuario con id " + event.ID)
		return err
	}

	logger.Info(fmt.Sprintf("Usuario eliminado: ID=%s", event.ID))

	return nil
}

func (h *UserEventHandler) HandleUserUpdated(data []byte) error {
	var user dto.UserRequest

	if err := json.Unmarshal(data, &user); err != nil {
		logger.Error(fmt.Sprintf("Error al deserializar evento de usuario actualizado: %v", err))
		return err
	}

	updatedUser, err := h.UserRepository.UpdateUser(user)

	if err != nil {
		logger.Error("No se ha podido actualizar el usuario con id " + user.ID)
		return err
	}

	logger.Info(fmt.Sprintf("Usuario actualizado: ID=%s, Email=%s, Nombre=%s, Rol=%s",
		updatedUser.ID, updatedUser.Email, updatedUser.Name, updatedUser.Rol))

	return nil
}
