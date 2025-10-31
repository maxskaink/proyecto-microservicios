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
)

// Event es la estructura base para todos los eventos
type Event struct {
	Type      EventType              `json:"type"`
	Timestamp time.Time              `json:"timestamp"`
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
