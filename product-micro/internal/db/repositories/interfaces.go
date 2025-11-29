package repositories

import "github.com/maxskaink/proyecto-microservicios/product-micro/internal/dto"

// UserRepository define operaciones de persistencia para usuarios.
type IProductRepository interface {
	GetByIdProduct(id string, tenantID string) (*dto.ProductDTOResponse, error)
	ListProducts(page int, pageSize int, tenantID string) (*[]dto.ProductDTOResponse, error)
	CreateProduct(product *dto.ProductDTORequest, tenantID string) (*dto.ProductDTOResponse, error)
	UpdateProduct(id string, product *dto.ProductDTORequest, tenantID string) (*dto.ProductDTOResponse, error)
	DeleteProduct(id string, tenantID string) error
}

type IUserRepository interface {
	CreateUser(dto.UserRequest, string) (*dto.UserResponse, error)
	FindById(id string, tenantID string) (*dto.UserResponse, error)
	FindByUID(uid string, tenantID string) (*dto.UserResponse, error)
	UpdateUser(dto.UserRequest, string) (*dto.UserResponse, error)
	DeleteUser(id string, tenantID string) error
}
