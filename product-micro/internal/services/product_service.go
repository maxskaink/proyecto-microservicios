package services

import (
	"fmt"

	"github.com/maxskaink/proyecto-microservicios/product-micro/internal/db/repositories"
	"github.com/maxskaink/proyecto-microservicios/product-micro/internal/domain"
	"github.com/maxskaink/proyecto-microservicios/product-micro/internal/dto"
)

type productService struct {
	productRepo repositories.IProductRepository
	userService IUserService
}

func NewProductService(productRepo repositories.IProductRepository, userService IUserService) IProductService {
	return &productService{
		productRepo: productRepo,
		userService: userService,
	}
}

// CreateProduct implements IProductService.
func (p *productService) CreateProduct(product dto.ProductDTORequest, idProducer string) (*dto.ProductDTOResponse, error) {
	// Validar que el usuario existe
	user, err := p.userService.GetUserByUUID(idProducer)
	if err != nil {
		return nil, domain.NotFoundError{Message: "Usuario con id " + idProducer + " no encontrado"}
	}

	// Validar que el usuario tiene rol de productor o administrador
	if user.Rol != domain.UserRoleProducer && user.Rol != domain.UserRoleAdmin {
		return nil, domain.UnauthorizedError{Message: "El usuario no tiene permisos para crear un producto"}
	}

	// Asignar el ProducerID al producto
	product.ProducerID = user.ID

	// Crear el producto en la base de datos
	result, err := p.productRepo.CreateProduct(&product)
	if err != nil {
		return nil, err
	}

	return result, nil
}

// GetByIdProduct implements IProductService.
func (p *productService) GetByIdProduct(id string) (*dto.ProductDTOResponse, error) {
	product, err := p.productRepo.GetByIdProduct(id)
	if err != nil {
		return nil, err
	}

	if product == nil {
		return nil, domain.NotFoundError{Message: fmt.Sprintf("El producto con id %s no se ha encontrado", id)}
	}

	return product, nil
}

// ListProduct implements IProductService.
func (p *productService) ListProduct(page int, pageSize int) (*[]dto.ProductDTOResponse, error) {

	//Validate the paramaters
	if page < 1 {
		return nil, domain.BadRequestError{Message: "La pagina no puede ser menor a 1"}
	}

	if pageSize < 2 || pageSize > 100 {
		return nil, domain.BadRequestError{Message: "El tamanio de la pagina debe ser entre 2 y 100"}
	}

	products, err := p.productRepo.ListProducts(1, 10)
	if err != nil {
		return nil, err
	}

	if products == nil {
		return nil, domain.NotFoundError{Message: "No se encontraron productos"}
	}

	return products, nil
}

// UpdateProduct implements IProductService.
func (p *productService) UpdateProduct(id string, product dto.ProductDTORequest, idProducer string) (*dto.ProductDTOResponse, error) {
	// Validar que el usuario existe
	user, err := p.userService.GetUserByUUID(idProducer)
	if err != nil {
		return nil, domain.NotFoundError{Message: fmt.Sprintf("Usuario con id %s no se ha encontrado", idProducer)}
	}

	// Obtener el producto actual
	currentProduct, err := p.productRepo.GetByIdProduct(id)
	if err != nil || currentProduct == nil {
		return nil, domain.NotFoundError{Message: fmt.Sprintf("Producto con id %s no fue encontrado", id)}
	}

	// Validar que el usuario es el dueño del producto o es administrador
	if currentProduct.ProducerID != idProducer && user.Rol != domain.UserRoleAdmin {
		return nil, domain.UnauthorizedError{Message: "No tiene permisos para actualizar el producto"}
	}

	product.ProducerID = user.ID

	result, err := p.productRepo.UpdateProduct(id, &product)

	if err != nil {
		return nil, err
	}

	return result, nil
}
