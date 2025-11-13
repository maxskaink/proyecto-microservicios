package handlers

import (
	"fmt"

	"github.com/maxskaink/proyecto-microservicios/product-micro/internal/db/repositories"
	"github.com/maxskaink/proyecto-microservicios/product-micro/internal/services/tenant"
	"github.com/maxskaink/proyecto-microservicios/product-micro/pkg/logger"
)

// EventDispatcher coordina los handlers de eventos
type EventDispatcher struct {
	registry *EventRegistry
}

// NewEventDispatcher crea un nuevo dispatcher con los handlers registrados
func NewEventDispatcher(
	userRepository repositories.IUserRepository,
	productRepository repositories.IProductRepository,
	tenantService *tenant.TenantService,
) *EventDispatcher {
	registry := NewEventRegistry()

	// Registrar handlers de usuario
	registry.Register(NewUserCreatedHandler(userRepository))
	registry.Register(NewUserDeletedHandler(userRepository))
	registry.Register(NewUserUpdatedHandler(userRepository))

	// Registrar handlers de tenant
	registry.Register(NewTenantCreatedHandler(tenantService))
	registry.Register(NewTenantDeletedHandler(tenantService))

	// Registrar handlers de orden
	registry.Register(NewOrderPaidHandler(productRepository))

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
	// El routingKey es el eventType (ej: "user.created", "tenant.deleted")
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
