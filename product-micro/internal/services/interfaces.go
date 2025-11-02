package services

import (
	"github.com/maxskaink/proyecto-microservicios/product-micro/internal/dto"
)

// IUserService define la lógica de negocio para usuarios.
type IUserService interface {
	GetUserByID(id string, tenantID string) (dto.UserResponse, error)
	GetUserByUUID(id string, tenantID string) (dto.UserResponse, error)
}

type IProductService interface {
	ListProduct(page int, pageSize int, tenantID string) (*[]dto.ProductDTOResponse, error)
	GetByIdProduct(id string, tenantID string) (*dto.ProductDTOResponse, error)
	CreateProduct(product dto.ProductDTORequest, idProducer string, tenantID string) (*dto.ProductDTOResponse, error)
	UpdateProduct(idProduct string, product dto.ProductDTORequest, idProducer string, tenantID string) (*dto.ProductDTOResponse, error)
}
