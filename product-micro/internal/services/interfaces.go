package services

import (
	"github.com/maxskaink/proyecto-microservicios/product-micro/internal/dto"
)

// IUserService define la lógica de negocio para usuarios.
type IUserService interface {
	GetUserByID(id string) (dto.UserResponse, error)
	GetUserByUUID(id string) (dto.UserResponse, error)
}

type IProductService interface {
	ListProduct(page int, pageSize int) (*[]dto.ProductDTOResponse, error)
	GetByIdProduct(id string) (*dto.ProductDTOResponse, error)
	CreateProduct(product dto.ProductDTORequest, idProducer string) (*dto.ProductDTOResponse, error)
	UpdateProduct(idProduct string, product dto.ProductDTORequest, idProducer string) (*dto.ProductDTOResponse, error)
}
