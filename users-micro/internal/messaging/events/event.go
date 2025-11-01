package events

import "time"

// EventType define los tipos de eventos
type EventType string

const (
	UserCreated   EventType = "user.created"
	UserUpdated   EventType = "user.updated"
	UserDeleted   EventType = "user.deleted"
	TenantCreated EventType = "tenant.created"
	TenantDeleted EventType = "tenant.deleted"
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
