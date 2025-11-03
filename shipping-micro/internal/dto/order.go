package dto

import (
	"time"
)

// OrderDTO representa una orden
type OrderDTO struct {
	ID         string         `json:"id"`
	UserID     string         `json:"user_id"`
	TotalPrice float64        `json:"total_price"`
	Status     string         `json:"status"`
	Items      []OrderItemDTO `json:"items,omitempty"`
	CreatedAt  time.Time      `json:"created_at"`
	UpdatedAt  time.Time      `json:"updated_at"`
}

// OrderItemDTO representa un producto en una orden
type OrderItemDTO struct {
	ID        string      `json:"id"`
	OrderID   string      `json:"order_id"`
	ProductID string      `json:"product_id"`
	Product   *ProductDTO `json:"product,omitempty"`
	Quantity  int         `json:"quantity"`
	Price     float64     `json:"price"`
	CreatedAt time.Time   `json:"created_at"`
}

// CreateOrderRequest solicitud para crear una orden desde el carrito
type CreateOrderRequest struct {
	ShippingAddress string `json:"shipping_address" binding:"required"`
}

// UpdateOrderStatusRequest solicitud para actualizar estado de orden
type UpdateOrderStatusRequest struct {
	Status string `json:"status" binding:"required,oneof=pending paid cancelled"`
}
