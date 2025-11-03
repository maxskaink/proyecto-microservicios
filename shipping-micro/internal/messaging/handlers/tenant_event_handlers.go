package handlers

import (
	"encoding/json"
	"fmt"

	"github.com/maxskaink/proyecto-microservicios/shipping-micro/internal/messaging/events"
	tenantSvc "github.com/maxskaink/proyecto-microservicios/shipping-micro/internal/services/tenant"
	"github.com/maxskaink/proyecto-microservicios/shipping-micro/pkg/logger"
)

// TenantCreatedHandler
type TenantCreatedHandler struct {
	BaseHandler
	tenant *tenantSvc.TenantService
}

func NewTenantCreatedHandler(s *tenantSvc.TenantService) *TenantCreatedHandler {
	return &TenantCreatedHandler{BaseHandler: NewBaseHandler(string(events.TenantCreated)), tenant: s}
}
func (h *TenantCreatedHandler) Handle(data []byte) error {
	// Soportar envelope y plano
	var env events.Event
	if err := json.Unmarshal(data, &env); err == nil && env.Data != nil {
		// event.data.tenant_id
		if tid, ok := env.Data["tenant_id"].(string); ok {
			return h.tenant.CreateTenantSchema(tid)
		}
		logger.Error("tenant_id no encontrado en data")
		return fmt.Errorf("tenant_id not found")
	}
	var m map[string]interface{}
	if err := json.Unmarshal(data, &m); err != nil {
		return err
	}
	if tid, ok := m["tenant_id"].(string); ok {
		fmt.Printf("Creando esquema para tenant_id: %s\n", tid)
		return h.tenant.CreateTenantSchema(tid)
	}
	return fmt.Errorf("tenant_id not found")
}

// TenantDeletedHandler
type TenantDeletedHandler struct {
	BaseHandler
	tenant *tenantSvc.TenantService
}

func NewTenantDeletedHandler(s *tenantSvc.TenantService) *TenantDeletedHandler {
	return &TenantDeletedHandler{BaseHandler: NewBaseHandler(string(events.TenantDeleted)), tenant: s}
}
func (h *TenantDeletedHandler) Handle(data []byte) error {
	var env events.Event
	if err := json.Unmarshal(data, &env); err == nil && env.Data != nil {
		if tid, ok := env.Data["tenant_id"].(string); ok {
			return h.tenant.DeleteTenantSchema(tid)
		}
		logger.Error("tenant_id no encontrado en data")
		return fmt.Errorf("tenant_id not found")
	}
	var m map[string]interface{}
	if err := json.Unmarshal(data, &m); err != nil {
		return err
	}
	if tid, ok := m["tenant_id"].(string); ok {
		return h.tenant.DeleteTenantSchema(tid)
	}
	return fmt.Errorf("tenant_id not found")
}
