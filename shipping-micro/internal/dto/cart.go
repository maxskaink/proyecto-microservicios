package dto

import (
	"time"
)

// CartItemDTO representa un artículo en el carrito
type CartItemDTO struct {
	ID        string      `json:"id"`
	UserID    string      `json:"user_id"`
	ProductID string      `json:"product_id"`
	Product   *ProductDTO `json:"product,omitempty"`
	Quantity  int         `json:"quantity"`
	CreatedAt time.Time   `json:"created_at"`
	UpdatedAt time.Time   `json:"updated_at"`
}

// AddToCartRequest solicitud para agregar producto al carrito
type AddToCartRequest struct {
	ProductID string `json:"product_id" binding:"required"`
	Quantity  int    `json:"quantity" binding:"required,min=1"`
}

// UpdateCartItemRequest solicitud para actualizar cantidad en carrito
type UpdateCartItemRequest struct {
	Quantity int `json:"quantity" binding:"required,min=0"`
}
