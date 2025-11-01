package events

import "time"

// EventType define los tipos de eventos
type EventType string

const (
	UserCreated    EventType = "user.created"
	UserUpdated    EventType = "user.updated"
	UserDeleted    EventType = "user.deleted"
	ProductCreated EventType = "product.created"
	ProductUpdated EventType = "product.updated"
	TenantCreated  EventType = "tenant.created"
	TenantDeleted  EventType = "tenant.deleted"
)

// Event es la estructura base para todos los eventos
type Event struct {
	Type      EventType              `json:"type"`
	Timestamp time.Time              `json:"timestamp"`
	TenantID  string                 `json:"tenant_id"`
	Data      map[string]interface{} `json:"data"`
}

// UserDeletedEvent representa un evento de usuario eliminado
type UserDeletedEvent struct {
	Event
	ID string `json:"id"`
}

// ProductCreatedEvent representa un evento de producto creado
type ProductCreatedEvent struct {
	Event
	ProductID  string `json:"product_id"`
	ProducerID string `json:"producer_id"`
}

// ProductUpdatedEvent representa un evento de producto actualizado
type ProductUpdatedEvent struct {
	Event
	ProductID  string `json:"product_id"`
	ProducerID string `json:"producer_id"`
}

// TenantCreatedEvent representa un evento de creación de tenant
type TenantCreatedEvent struct {
	Event
	TenantID   string `json:"tenant_id"`
	TenantName string `json:"tenant_name"`
}

// TenantDeletedEvent representa un evento de eliminación de tenant
type TenantDeletedEvent struct {
	Event
	TenantID string `json:"tenant_id"`
}
