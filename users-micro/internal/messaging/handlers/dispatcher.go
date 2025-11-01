package handlers

import (
	"context"
	"fmt"

	"github.com/maxskaink/proyecto-microservicios/users-micro/internal/services/tenant"
	"github.com/maxskaink/proyecto-microservicios/users-micro/pkg/logger"
)

// EventDispatcher coordina los handlers de eventos
type EventDispatcher struct {
	registry *EventRegistry
}

// NewEventDispatcher crea un nuevo dispatcher con los handlers registrados
func NewEventDispatcher(tenantService *tenant.TenantService) *EventDispatcher {
	registry := NewEventRegistry()

	// Registrar handlers de tenant
	registry.Register(NewTenantCreatedHandler(tenantService))
	registry.Register(NewTenantDeletedHandler(tenantService))

	logger.Info("EventDispatcher inicializado con todos los handlers")

	return &EventDispatcher{
		registry: registry,
	}
}

// Dispatch procesa un evento basado en su tipo (compatible con routing key)
func (ed *EventDispatcher) Dispatch(eventType string, data []byte) error {
	return ed.registry.Handle(eventType, data)
}

// ProcessMessage implementa la interfaz Dispatcher de rabbitmq
// Recibe el mensaje del consumer y lo rutea al handler apropiado
func (ed *EventDispatcher) ProcessMessage(data []byte, routingKey string) error {
	// El routingKey es el eventType (ej: "tenant.created", "tenant.deleted")
	return ed.registry.Handle(routingKey, data)
}

// RegisterHandler registra un nuevo handler dinámicamente
func (ed *EventDispatcher) RegisterHandler(handler EventHandler) error {
	return ed.registry.Register(handler)
}

// ListHandlers retorna todos los tipos de eventos registrados
func (ed *EventDispatcher) ListHandlers() []string {
	return ed.registry.ListHandlers()
}

// Info retorna información del dispatcher
func (ed *EventDispatcher) Info() string {
	handlers := ed.ListHandlers()
	return fmt.Sprintf("EventDispatcher con %d handlers registrados: %v", len(handlers), handlers)
}

// Close cierra el dispatcher (no hay recursos que cerrar)
func (ed *EventDispatcher) Close() error {
	return nil
}

// StartConsuming inicia el consumo de eventos (implementa la interfaz Consumer)
func (ed *EventDispatcher) StartConsuming(ctx context.Context) error {
	// Este método es implementado por el Consumer, no por el Dispatcher
	// El Dispatcher solo procesa mensajes
	return nil
}
