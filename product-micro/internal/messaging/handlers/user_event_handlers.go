package handlers

import (
	"encoding/json"
	"fmt"

	"github.com/maxskaink/proyecto-microservicios/product-micro/internal/db/repositories"
	"github.com/maxskaink/proyecto-microservicios/product-micro/internal/dto"
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
	var event map[string]interface{}
	if err := h.UnmarshalEvent(data, &event); err != nil {
		return err
	}

	// Extraer payload
	payload, ok := event["payload"]
	if !ok {
		// Si no hay payload, podría ser un evento de tenant que llegó por error
		// Loguear advertencia pero no fallar
		logger.Error(fmt.Sprintf("No se encontró 'payload' en evento user.created. Event: %v", event))
		return fmt.Errorf("payload not found in user.created event - possible routing issue")
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

	// Extraer tenant_id del evento
	tenantID, ok := event["tenant_id"].(string)
	if !ok {
		logger.Error("tenant_id no encontrado o no es string")
		return fmt.Errorf("tenant_id not found")
	}

	// Crear usuario
	userCreated, err := h.userRepo.CreateUser(user, tenantID)
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
	var event map[string]interface{}
	if err := h.UnmarshalEvent(data, &event); err != nil {
		return err
	}

	// Extraer campos del evento
	userID, ok := event["user_id"].(string)
	if !ok {
		logger.Error("user_id no encontrado o no es string")
		return fmt.Errorf("user_id not found")
	}

	tenantID, ok := event["tenant_id"].(string)
	if !ok {
		logger.Error("tenant_id no encontrado o no es string")
		return fmt.Errorf("tenant_id not found")
	}

	// Eliminar usuario
	if err := h.userRepo.DeleteUser(userID, tenantID); err != nil {
		logger.Error(fmt.Sprintf("No se pudo eliminar el usuario con id %s: %v", userID, err))
		return err
	}

	h.LogEvent(fmt.Sprintf("Usuario eliminado: ID=%s", userID))
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
	var event map[string]interface{}
	if err := h.UnmarshalEvent(data, &event); err != nil {
		return err
	}

	// Extraer payload
	payload, ok := event["payload"]
	if !ok {
		logger.Error("No se encontró 'payload' en el evento")
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

	// Extraer tenant_id del evento
	tenantID, ok := event["tenant_id"].(string)
	if !ok {
		logger.Error("tenant_id no encontrado o no es string")
		return fmt.Errorf("tenant_id not found")
	}

	// Aquí iría la lógica de actualización del usuario
	// TODO: Implementar actualización en repositorio
	_ = tenantID // usar tenantID cuando se implemente

	h.LogEvent(fmt.Sprintf("Usuario actualizado: ID=%s", user.ID))
	return nil
}
