package handlers

import (
	"github.com/maxskaink/proyecto-microservicios/shipping-micro/internal/db/repositories"
	tenant "github.com/maxskaink/proyecto-microservicios/shipping-micro/internal/services/tenant"
)

// Dispatcher que adapta routingKey a handler por tipo
type EventDispatcher struct{ registry *EventRegistry }

func NewEventDispatcher(productRepo repositories.IProductRepository, userRepo repositories.IUserRepository, tenantService *tenant.TenantService) *EventDispatcher {
	reg := NewEventRegistry()
	// Registrar handlers
	reg.Register(NewProductCreatedHandler(productRepo))
	reg.Register(NewProductUpdatedHandler(productRepo))
	reg.Register(NewUserCreatedHandler(userRepo))
	reg.Register(NewTenantCreatedHandler(tenantService))
	reg.Register(NewTenantDeletedHandler(tenantService))
	return &EventDispatcher{registry: reg}
}

func (ed *EventDispatcher) Dispatch(eventType string, data []byte) error {
	return ed.registry.Handle(eventType, data)
}

func (ed *EventDispatcher) ProcessMessage(data []byte, routingKey string) error {
	return ed.registry.Handle(routingKey, data)
}

func (ed *EventDispatcher) RegisterHandler(handler EventHandler) error {
	return ed.registry.Register(handler)
}

func (ed *EventDispatcher) ListHandlers() []string { return ed.registry.ListHandlers() }
