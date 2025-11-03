package dto

import (
	"time"
)

// ShippingDTO representa un envío
type ShippingDTO struct {
	ID              string    `json:"id"`
	OrderID         string    `json:"order_id"`
	TrackingNumber  string    `json:"tracking_number"`
	ShippingAddress string    `json:"shipping_address"`
	Status          string    `json:"status"`
	CreatedAt       time.Time `json:"created_at"`
	UpdatedAt       time.Time `json:"updated_at"`
}

// UpdateShippingStatusRequest solicitud para actualizar estado de envío
type UpdateShippingStatusRequest struct {
	Status string `json:"status" binding:"required,oneof=pending in_transit delivered cancelled"`
}
