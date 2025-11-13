package handlers

import (
	"encoding/json"
	"fmt"

	"github.com/maxskaink/proyecto-microservicios/shipping-micro/internal/db/repositories"
	"github.com/maxskaink/proyecto-microservicios/shipping-micro/internal/dto"
	"github.com/maxskaink/proyecto-microservicios/shipping-micro/internal/messaging/events"
	"github.com/maxskaink/proyecto-microservicios/shipping-micro/pkg/logger"
)

// UserCreatedHandler maneja user.created
type UserCreatedHandler struct {
	BaseHandler
	repo repositories.IUserRepository
}

func NewUserCreatedHandler(repo repositories.IUserRepository) *UserCreatedHandler {
	return &UserCreatedHandler{BaseHandler: NewBaseHandler(string(events.UserCreated)), repo: repo}
}
func (h *UserCreatedHandler) Handle(data []byte) error {
	fmt.Printf("Creando usuario por evento...\n")
	// Formato envelope
	var env events.Event
	if err := json.Unmarshal(data, &env); err == nil && env.Data != nil {
		payload, ok := env.Data["payload"]
		if !ok {
			logger.Error("payload no encontrado en evento user.created")
			return fmt.Errorf("payload not found")
		}
		b, _ := json.Marshal(payload)
		var u dto.UserDTO
		if err := json.Unmarshal(b, &u); err != nil {
			return err
		}
		fmt.Println("Creating user:")
		fmt.Println(u)
		return h.repo.Create(&u, env.TenantID)
	}
	// Formato plano
	var u dto.UserDTO
	if err := json.Unmarshal(data, &u); err != nil {
		return err
	}
	return h.repo.Create(&u, "")
}

// UserUpdatedHandler maneja user.updated
type UserUpdatedHandler struct {
	BaseHandler
	repo repositories.IUserRepository
}

func NewUserUpdatedHandler(repo repositories.IUserRepository) *UserUpdatedHandler {
	return &UserUpdatedHandler{BaseHandler: NewBaseHandler(string(events.UserUpdated)), repo: repo}
}

func (h *UserUpdatedHandler) Handle(data []byte) error {
	// Soporta envelope y formato plano
	var env events.Event
	if err := json.Unmarshal(data, &env); err == nil && env.Data != nil {
		payload, ok := env.Data["payload"]
		if !ok {
			logger.Error("payload no encontrado en evento user.updated")
			return fmt.Errorf("payload not found")
		}
		b, _ := json.Marshal(payload)
		var u dto.UserDTO
		if err := json.Unmarshal(b, &u); err != nil {
			return err
		}
		return h.repo.Update(&u, env.TenantID)
	}
	// Formato plano
	var u dto.UserDTO
	if err := json.Unmarshal(data, &u); err != nil {
		return err
	}
	return h.repo.Update(&u, "")
}
