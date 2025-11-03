package events

import "time"

type EventType string

const (
	UserCreated        EventType = "user.created"
	UserUpdated        EventType = "user.updated"
	UserDeleted        EventType = "user.deleted"
	ProductCreated     EventType = "product.created"
	ProductUpdated     EventType = "product.updated"
	TenantCreated      EventType = "tenant.created"
	TenantDeleted      EventType = "tenant.deleted"
	OrderCreated       EventType = "order.created"
	OrderStatusChanged EventType = "order.status_changed"
	ShippingCreated    EventType = "shipping.created"
)

// Envelope compatible con otros micros (users-micro)
type Event struct {
	Type      EventType              `json:"type"`
	Timestamp time.Time              `json:"timestamp"`
	TenantID  string                 `json:"tenant_id"`
	Data      map[string]interface{} `json:"data"`
}
