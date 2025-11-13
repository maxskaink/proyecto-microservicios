package dto

import (
	"time"
)

// ProductDTO representa un producto
type ProductDTO struct {
	ID          string    `json:"id"`
	Name        string    `json:"name"`
	ProducerID  string    `json:"producer_id"`
	Description string    `json:"description"`
	Price       float64   `json:"price"`
	Stock       int       `json:"stock"`
	CreatedAt   time.Time `json:"created_at"`
	UpdatedAt   time.Time `json:"updated_at"`
}
