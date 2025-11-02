package handlers

import (
	"context"
	"fmt"

	"github.com/maxskaink/proyecto-microservicios/users-micro/internal/services/tenant"
	"github.com/maxskaink/proyecto-microservicios/users-micro/pkg/logger"
)

// TenantCreatedHandler maneja eventos de creación de tenant
type TenantCreatedHandler struct {
	BaseHandler
	tenantService *tenant.TenantService
}

// NewTenantCreatedHandler crea un nuevo handler de creación de tenant
func NewTenantCreatedHandler(ts *tenant.TenantService) *TenantCreatedHandler {
	return &TenantCreatedHandler{
		BaseHandler:   NewBaseHandler("tenant.created"),
		tenantService: ts,
	}
}

// Handle procesa el evento de creación de tenant
func (h *TenantCreatedHandler) Handle(data []byte) error {
	var event map[string]interface{}
	if err := h.UnmarshalEvent(data, &event); err != nil {
		return err
	}

	// Extraer campos del evento
	eventData, ok := event["data"].(map[string]interface{})
	if !ok {
		logger.Error("'data' no encontrado en evento de tenant creado")
		return fmt.Errorf("data not found")
	}

	tenantID, ok := eventData["tenant_id"].(string)
	if !ok {
		logger.Error("tenant_id no encontrado o no es string")
		return fmt.Errorf("tenant_id not found")
	}

	tenantName, ok := eventData["tenant_name"].(string)
	if !ok {
		logger.Error("tenant_name no encontrado o no es string")
		return fmt.Errorf("tenant_name not found")
	}

	// Crear el schema para el tenant
	if err := h.tenantService.CreateTenant(context.Background(), tenantID, tenantName); err != nil {
		logger.Error(fmt.Sprintf("Error al crear tenant %s: %v", tenantID, err))
		return err
	}

	h.LogEvent(fmt.Sprintf("Tenant creado: ID=%s, Name=%s", tenantID, tenantName))
	return nil
}

// TenantDeletedHandler maneja eventos de eliminación de tenant
type TenantDeletedHandler struct {
	BaseHandler
	tenantService *tenant.TenantService
}

// NewTenantDeletedHandler crea un nuevo handler de eliminación de tenant
func NewTenantDeletedHandler(ts *tenant.TenantService) *TenantDeletedHandler {
	return &TenantDeletedHandler{
		BaseHandler:   NewBaseHandler("tenant.deleted"),
		tenantService: ts,
	}
}

// Handle procesa el evento de eliminación de tenant
func (h *TenantDeletedHandler) Handle(data []byte) error {
	var event map[string]interface{}
	if err := h.UnmarshalEvent(data, &event); err != nil {
		return err
	}

	// Extraer campos del evento
	eventData, ok := event["data"].(map[string]interface{})
	if !ok {
		logger.Error("'data' no encontrado en evento de tenant eliminado")
		return fmt.Errorf("data not found")
	}

	tenantID, ok := eventData["tenant_id"].(string)
	if !ok {
		logger.Error("tenant_id no encontrado o no es string")
		return fmt.Errorf("tenant_id not found")
	}

	// Eliminar el schema del tenant
	if err := h.tenantService.DeleteTenant(context.Background(), tenantID); err != nil {
		logger.Error(fmt.Sprintf("Error al eliminar tenant %s: %v", tenantID, err))
		return err
	}

	h.LogEvent(fmt.Sprintf("Tenant eliminado: ID=%s", tenantID))
	return nil
}
