package domain

import (
	"time"

	"github.com/google/uuid"
)

// Product representa un producto en el dominio
type Product struct {
	ID          uuid.UUID
	Name        string
	Description string
	Price       float64
	Stock       int
	CreatedAt   time.Time
	UpdatedAt   time.Time
}

// User representa un usuario en el dominio
type User struct {
	ID        uuid.UUID
	Email     string
	Name      string
	CreatedAt time.Time
	UpdatedAt time.Time
}

// CartItem representa un artículo en el carrito
type CartItem struct {
	ID        uuid.UUID
	UserID    uuid.UUID
	ProductID uuid.UUID
	Quantity  int
	CreatedAt time.Time
	UpdatedAt time.Time
}

// OrderStatus representa el estado de una orden
type OrderStatus string

const (
	OrderStatusPending   OrderStatus = "pending"
	OrderStatusPaid      OrderStatus = "paid"
	OrderStatusCancelled OrderStatus = "cancelled"
)

var orderStatus = map[OrderStatus]struct{}{
	OrderStatusPending:   {},
	OrderStatusPaid:      {},
	OrderStatusCancelled: {},
}

// Order representa una orden/pedido
type Order struct {
	ID         uuid.UUID
	UserID     uuid.UUID
	TotalPrice float64
	Status     OrderStatus
	CreatedAt  time.Time
	UpdatedAt  time.Time
}

// OrderItem representa un producto dentro de una orden
type OrderItem struct {
	ID        uuid.UUID
	OrderID   uuid.UUID
	ProductID uuid.UUID
	Quantity  int
	Price     float64
	CreatedAt time.Time
}

// ShippingStatus representa el estado de un envío
type ShippingStatus string

const (
	ShippingStatusPending   ShippingStatus = "pending"
	ShippingStatusInTransit ShippingStatus = "in_transit"
	ShippingStatusDelivered ShippingStatus = "delivered"
	ShippingStatusCancelled ShippingStatus = "cancelled"
)

var shippingStatus = map[ShippingStatus]struct{}{
	ShippingStatusPending:   {},
	ShippingStatusInTransit: {},
	ShippingStatusDelivered: {},
	ShippingStatusCancelled: {},
}

// Shipping representa un envío
type Shipping struct {
	ID              uuid.UUID
	OrderID         uuid.UUID
	TrackingNumber  string
	ShippingAddress string
	Status          ShippingStatus
	CreatedAt       time.Time
	UpdatedAt       time.Time
}

func IsValidOrderStatus(status string) bool {
	_, exists := orderStatus[OrderStatus(status)]
	return exists
}

func IsValidShippingStatus(status string) bool {
	_, exists := shippingStatus[ShippingStatus(status)]
	return exists
}
