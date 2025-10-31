package repositories

import "github.com/maxskaink/proyecto-microservicios/product-micro/internal/dto"

// UserRepository define operaciones de persistencia para usuarios.
type IProductRepository interface {
	GetByIdProduct(id string) (*dto.ProductDTOResponse, error)
	ListProducts(page int, pageSize int) (*[]dto.ProductDTOResponse, error)
	CreateProduct(product *dto.ProductDTORequest) (*dto.ProductDTOResponse, error)
	UpdateProduct(id string, product *dto.ProductDTORequest) (*dto.ProductDTOResponse, error)
}

type IUserRepository interface {
	CreateUser(dto.UserRequest) (*dto.UserResponse, error)
	FindById(id string) (*dto.UserResponse, error)
	FindByUID(uid string) (*dto.UserResponse, error)
}
