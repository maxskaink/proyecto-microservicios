package handlers

import (
	"encoding/json"
	"fmt"

	"github.com/maxskaink/proyecto-microservicios/product-micro/internal/db/repositories"
	"github.com/maxskaink/proyecto-microservicios/product-micro/internal/dto"
	"github.com/maxskaink/proyecto-microservicios/product-micro/internal/messaging/events"
	"github.com/maxskaink/proyecto-microservicios/product-micro/pkg/logger"
)

// UserCreatedHandler maneja eventos de creación de usuario
type UserCreatedHandler struct {
	BaseHandler
	userRepo repositories.IUserRepository
}

// NewUserCreatedHandler crea un nuevo handler de creación de usuario
func NewUserCreatedHandler(userRepo repositories.IUserRepository) *UserCreatedHandler {
	return &UserCreatedHandler{
		BaseHandler: NewBaseHandler("user.created"),
		userRepo:    userRepo,
	}
}

// Handle procesa el evento de creación de usuario
func (h *UserCreatedHandler) Handle(data []byte) error {
	var event events.Event
	if err := json.Unmarshal(data, &event); err != nil {
		logger.Error(fmt.Sprintf("Error al deserializar evento: %v", err))
		return err
	}

	// Extraer payload desde data
	payload, ok := event.Data["payload"]
	if !ok {
		logger.Error(fmt.Sprintf("No se encontró 'payload' en evento user.created. Event: %v", event))
		return fmt.Errorf("payload not found in user.created event")
	}

	// Serializar y deserializar payload a UserRequest
	userJSON, err := json.Marshal(payload)
	if err != nil {
		logger.Error(fmt.Sprintf("Error al serializar payload: %v", err))
		return err
	}

	var user dto.UserRequest
	if err := json.Unmarshal(userJSON, &user); err != nil {
		logger.Error(fmt.Sprintf("Error al deserializar a UserRequest: %v", err))
		return err
	}

	// Crear usuario
	userCreated, err := h.userRepo.CreateUser(user, event.TenantID)
	if err != nil {
		logger.Error(fmt.Sprintf("No se pudo crear el usuario con id %s: %v", user.ID, err))
		return err
	}

	h.LogEvent(fmt.Sprintf("Usuario creado: ID=%s, Email=%s, Nombre=%s", userCreated.ID, userCreated.Email, userCreated.Name))
	return nil
}

// UserDeletedHandler maneja eventos de eliminación de usuario
type UserDeletedHandler struct {
	BaseHandler
	userRepo repositories.IUserRepository
}

// NewUserDeletedHandler crea un nuevo handler de eliminación de usuario
func NewUserDeletedHandler(userRepo repositories.IUserRepository) *UserDeletedHandler {
	return &UserDeletedHandler{
		BaseHandler: NewBaseHandler("user.deleted"),
		userRepo:    userRepo,
	}
}

// Handle procesa el evento de eliminación de usuario
func (h *UserDeletedHandler) Handle(data []byte) error {
	var event events.UserDeletedEvent
	if err := json.Unmarshal(data, &event); err != nil {
		logger.Error(fmt.Sprintf("Error al deserializar evento: %v", err))
		return err
	}

	// Eliminar usuario
	if err := h.userRepo.DeleteUser(event.ID, event.TenantID); err != nil {
		logger.Error(fmt.Sprintf("No se pudo eliminar el usuario con id %s: %v", event.ID, err))
		return err
	}

	h.LogEvent(fmt.Sprintf("Usuario eliminado: ID=%s", event.ID))
	return nil
}

// UserUpdatedHandler maneja eventos de actualización de usuario
type UserUpdatedHandler struct {
	BaseHandler
	userRepo repositories.IUserRepository
}

// NewUserUpdatedHandler crea un nuevo handler de actualización de usuario
func NewUserUpdatedHandler(userRepo repositories.IUserRepository) *UserUpdatedHandler {
	return &UserUpdatedHandler{
		BaseHandler: NewBaseHandler("user.updated"),
		userRepo:    userRepo,
	}
}

// Handle procesa el evento de actualización de usuario
func (h *UserUpdatedHandler) Handle(data []byte) error {
	var event events.Event
	if err := json.Unmarshal(data, &event); err != nil {
		logger.Error(fmt.Sprintf("Error al deserializar evento: %v", err))
		return err
	}

	// Extraer payload desde data
	payload, ok := event.Data["payload"]
	if !ok {
		logger.Error("No se encontró 'payload' en data")
		return fmt.Errorf("payload not found")
	}

	// Serializar y deserializar payload a UserRequest
	userJSON, err := json.Marshal(payload)
	if err != nil {
		logger.Error(fmt.Sprintf("Error al serializar payload: %v", err))
		return err
	}

	var user dto.UserRequest
	if err := json.Unmarshal(userJSON, &user); err != nil {
		logger.Error(fmt.Sprintf("Error al deserializar a UserRequest: %v", err))
		return err
	}

	// Actualizar el usuario
	_, err = h.userRepo.UpdateUser(user, event.TenantID)
	if err != nil {
		logger.Error(fmt.Sprintf("No se pudo actualizar el usuario con id %s: %v", user.ID, err))
		return err
	}

	h.LogEvent(fmt.Sprintf("Usuario actualizado: ID=%s", user.ID))
	return nil
}
