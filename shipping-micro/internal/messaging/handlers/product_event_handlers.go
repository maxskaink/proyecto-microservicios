package handlers

import (
	"encoding/json"
	"fmt"

	"github.com/maxskaink/proyecto-microservicios/shipping-micro/internal/db/repositories"
	"github.com/maxskaink/proyecto-microservicios/shipping-micro/internal/dto"
	"github.com/maxskaink/proyecto-microservicios/shipping-micro/internal/messaging/events"
	"github.com/maxskaink/proyecto-microservicios/shipping-micro/pkg/logger"
)

// ProductCreatedHandler maneja product.created
type ProductCreatedHandler struct {
	BaseHandler
	repo repositories.IProductRepository
}

func NewProductCreatedHandler(repo repositories.IProductRepository) *ProductCreatedHandler {
	return &ProductCreatedHandler{BaseHandler: NewBaseHandler(string(events.ProductCreated)), repo: repo}
}
func (h *ProductCreatedHandler) Handle(data []byte) error {
	// Soporta dos formatos: envelope (Event) con data.payload, y plano
	var env events.Event
	if err := json.Unmarshal(data, &env); err == nil && env.Data != nil {
		payload, ok := env.Data["payload"]
		if !ok {
			logger.Error("payload no encontrado en data")
			return fmt.Errorf("payload not found")
		}
		b, _ := json.Marshal(payload)
		var p dto.ProductDTO
		if err := json.Unmarshal(b, &p); err != nil {
			return err
		}
		fmt.Println(p)

		return h.repo.Create(&p, env.TenantID)
	}
	return fmt.Errorf("invalid format of the event")
}

// ProductUpdatedHandler maneja product.updated
type ProductUpdatedHandler struct {
	BaseHandler
	repo repositories.IProductRepository
}

func NewProductUpdatedHandler(repo repositories.IProductRepository) *ProductUpdatedHandler {
	return &ProductUpdatedHandler{BaseHandler: NewBaseHandler(string(events.ProductUpdated)), repo: repo}
}
func (h *ProductUpdatedHandler) Handle(data []byte) error {
	var env events.Event
	if err := json.Unmarshal(data, &env); err == nil && env.Data != nil {
		payload, ok := env.Data["payload"]
		if !ok {
			logger.Error("payload no encontrado en data")
			return fmt.Errorf("payload not found")
		}
		b, _ := json.Marshal(payload)
		var p dto.ProductDTO
		if err := json.Unmarshal(b, &p); err != nil {
			return err
		}
		return h.repo.Update(&p, env.TenantID)
	}
	var p dto.ProductDTO
	if err := json.Unmarshal(data, &p); err != nil {
		return err
	}
	return h.repo.Update(&p, "")
}
