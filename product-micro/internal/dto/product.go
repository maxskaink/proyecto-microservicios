package dto

import (
	"time"

	"github.com/maxskaink/proyecto-microservicios/product-micro/internal/domain"
)

type ProductDTOResponse struct {
	ID          string                 `json:"id"`
	ProducerID  string                 `json:"producer_id"`
	Category    domain.ProductCategory `json:"category"`
	Price       int                    `json:"price"`
	Name        string                 `json:"name"`
	Description string                 `json:"description"`
	Stock       int                    `json:"stock"`
	Unit        domain.ProductUnit     `json:"unit"`
	PhotoUrl    string                 `json:"photo_url"`
	CreatedAt   time.Time              `json:"create_at"`
}

type ProductDTORequest struct {
	ProducerID  string                 `json:"-"`
	Category    domain.ProductCategory `json:"category" binding:"required,category"`
	Price       int                    `json:"price" binding:"required,min=1"`
	Description string                 `json:"description" binding:"required,min=10,max=255"`
	Name        string                 `json:"name" binding:"required,min=3,max=100"`
	Stock       int                    `json:"stock" binding:"required,min=1"`
	Unit        domain.ProductUnit     `json:"unit" binding:"required,unit"`
	PhotoUrl    string                 `json:"photo_url" binding:"omitempty,url"`
}

type ProductPhotoInfoDTO struct {
	UploadURL string `json:"upload_url"`
	ObjectKey string `json:"object_key"`
	PublicURL string `json:"public_url"`
	ExpiresIn int    `json:"expires_in"` // segundos
}

type UploadURLRequest struct {
	Filename    string `json:"filename" binding:"required"`
	ContentType string `json:"content_type" binding:"required"`
}

type CompletePhotoRequest struct {
	ObjectKey string `json:"object_key" binding:"required"`
}

func ProductDTOResponseTORequest(response ProductDTOResponse) ProductDTORequest {
	return ProductDTORequest{
		ProducerID:  response.ProducerID,
		Category:    response.Category,
		Price:       response.Price,
		Description: response.Description,
		Name:        response.Name,
		Stock:       response.Stock,
		Unit:        response.Unit,
		PhotoUrl:    response.PhotoUrl,
	}
}
