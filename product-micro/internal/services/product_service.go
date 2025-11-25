package services

import (
	"context"
	"fmt"
	"path/filepath"
	"strings"
	"time"

	"github.com/google/uuid"
	"github.com/maxskaink/proyecto-microservicios/product-micro/internal/config"
	"github.com/maxskaink/proyecto-microservicios/product-micro/internal/db/repositories"
	"github.com/maxskaink/proyecto-microservicios/product-micro/internal/domain"
	"github.com/maxskaink/proyecto-microservicios/product-micro/internal/dto"
	"github.com/maxskaink/proyecto-microservicios/product-micro/internal/messaging"
	"github.com/maxskaink/proyecto-microservicios/product-micro/internal/storage"
)

type productService struct {
	productRepo repositories.IProductRepository
	userService IUserService
	publisher   messaging.Publisher

	//for storage
	storage storage.ObjectStorage
	cfg     *config.Config
}

func NewProductService(productRepo repositories.IProductRepository, userService IUserService, publisher messaging.Publisher, storage storage.ObjectStorage, cfg *config.Config) IProductService {
	return &productService{
		productRepo: productRepo,
		userService: userService,
		publisher:   publisher,
		storage:     storage,
		cfg:         cfg,
	}
}

// CreateProduct implements IProductService.
func (p *productService) CreateProduct(product dto.ProductDTORequest, idProducer string, tenantID string) (*dto.ProductDTOResponse, error) {
	// Validar tenant
	if tenantID == "" {
		return nil, domain.BadRequestError{Message: "Tenant requerido"}
	}

	// Validar que el usuario existe
	user, err := p.userService.GetUserByUUID(idProducer, tenantID)
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
	result, err := p.productRepo.CreateProduct(&product, tenantID)
	if err != nil {
		return nil, err
	}

	// Publicar evento de producto creado
	if p.publisher != nil && result != nil {
		// Ignoramos el error de publicación para no bloquear el flujo principal
		_ = p.publisher.PublishProductCreated(*result, tenantID)
	}

	return result, nil
}

// GetByIdProduct implements IProductService.
func (p *productService) GetByIdProduct(id string, tenantID string) (*dto.ProductDTOResponse, error) {
	// Validar tenant
	if tenantID == "" {
		return nil, domain.BadRequestError{Message: "Tenant requerido"}
	}

	product, err := p.productRepo.GetByIdProduct(id, tenantID)
	if err != nil {
		return nil, err
	}

	if product == nil {
		return nil, domain.NotFoundError{Message: fmt.Sprintf("El producto con id %s no se ha encontrado", id)}
	}

	return product, nil
}

// ListProduct implements IProductService.
func (p *productService) ListProduct(page int, pageSize int, tenantID string) (*[]dto.ProductDTOResponse, error) {
	// Validar tenant
	if tenantID == "" {
		return nil, domain.BadRequestError{Message: "Tenant requerido"}
	}

	//Validate the paramaters
	if page < 1 {
		return nil, domain.BadRequestError{Message: "La pagina no puede ser menor a 1"}
	}

	if pageSize < 2 || pageSize > 100 {
		return nil, domain.BadRequestError{Message: "El tamanio de la pagina debe ser entre 2 y 100"}
	}

	products, err := p.productRepo.ListProducts(page, pageSize, tenantID)
	if err != nil {
		return nil, err
	}

	if products == nil {
		return nil, domain.NotFoundError{Message: "No se encontraron productos"}
	}

	return products, nil
}

// UpdateProduct implements IProductService.
func (p *productService) UpdateProduct(id string, product dto.ProductDTORequest, idProducer string, tenantID string) (*dto.ProductDTOResponse, error) {
	// Validar tenant
	if tenantID == "" {
		return nil, domain.BadRequestError{Message: "Tenant requerido"}
	}

	// Validar que el usuario existe
	user, err := p.userService.GetUserByUUID(idProducer, tenantID)
	if err != nil {
		return nil, domain.NotFoundError{Message: fmt.Sprintf("Usuario con id %s no se ha encontrado", idProducer)}
	}

	// Obtener el producto actual
	currentProduct, err := p.productRepo.GetByIdProduct(id, tenantID)
	if err != nil || currentProduct == nil {
		return nil, domain.NotFoundError{Message: fmt.Sprintf("Producto con id %s no fue encontrado", id)}
	}

	// Validar que el usuario es el dueño del producto o es administrador
	if currentProduct.ProducerID != idProducer && user.Rol != domain.UserRoleAdmin {
		return nil, domain.UnauthorizedError{Message: "No tiene permisos para actualizar el producto"}
	}

	product.ProducerID = user.ID

	result, err := p.productRepo.UpdateProduct(id, &product, tenantID)

	if err != nil {
		return nil, err
	}

	// Publicar evento de producto actualizado
	if p.publisher != nil && result != nil {
		_ = p.publisher.PublishProductUpdated(*result, tenantID)

		// Si cambió el stock, publicar evento específico de actualización de stock
		if currentProduct.Stock != result.Stock {
			_ = p.publisher.PublishProductStockUpdated(*result, tenantID)
		}
	}

	return result, nil
}

// CompletePhotoUpload implements IProductService.
func (p *productService) CompletePhotoUpload(productID string, objectKey string, tenantID string, userUID string) (*dto.ProductDTOResponse, error) {
	if tenantID == "" {
		return nil, domain.BadRequestError{Message: "Tenant requerido"}
	}
	if productID == "" || objectKey == "" {
		return nil, domain.BadRequestError{Message: "Datos incompletos"}
	}
	if !strings.HasPrefix(objectKey, tenantID+"/") {
		return nil, domain.BadRequestError{Message: "object_key inválido"}
	}
	// Validar usuario
	user, err := p.userService.GetUserByUUID(userUID, tenantID)
	if err != nil {
		return nil, domain.UnauthorizedError{Message: "Usuario no válido"}
	}

	// Verificar producto
	current, err := p.productRepo.GetByIdProduct(productID, tenantID)
	if err != nil || current == nil {
		return nil, domain.NotFoundError{Message: "Producto no encontrado"}
	}
	if current.ProducerID != user.ID && user.Rol != domain.UserRoleAdmin {
		return nil, domain.UnauthorizedError{Message: "No autorizado a modificar foto"}
	}

	// HEAD objeto
	_, err = p.storage.HeadObject(context.Background(), p.cfg.R2Bucket, objectKey)
	if err != nil {
		return nil, domain.BadRequestError{Message: "Objeto no encontrado en storage"}
	}

	// Mover a final
	finalKey := fmt.Sprintf("%s/products/%s/%s", tenantID, productID, filepath.Base(objectKey))
	if err := p.storage.CopyObject(context.Background(), p.cfg.R2Bucket, objectKey, finalKey); err != nil {
		return nil, domain.InternalServerError{Message: "Error moviendo objeto"}
	}
	_ = p.storage.DeleteObject(context.Background(), p.cfg.R2Bucket, objectKey)

	publicURL := fmt.Sprintf("%s/%s", p.cfg.R2PublicBaseURL, finalKey)

	newProduct := dto.ProductDTOResponseTORequest(*current)
	newProduct.PhotoUrl = publicURL

	updated, err := p.productRepo.UpdateProduct(productID, &newProduct, tenantID)
	if err != nil {
		return nil, err
	}

	if p.publisher != nil {
		_ = p.publisher.PublishProductUpdated(*updated, tenantID)
	}

	return updated, nil
}

// GetUploadURL implements IProductService.
func (p *productService) GetUploadURL(tenantID string, filename string, contentType string) (*dto.ProductPhotoInfoDTO, error) {
	if tenantID == "" {
		return nil, domain.BadRequestError{Message: "Tenant requerido"}
	}
	// Validar content-type simple
	if !strings.HasPrefix(contentType, "image/") {
		return nil, domain.BadRequestError{Message: "Tipo de contenido no permitido"}
	}
	ext := filepath.Ext(filename)
	if ext == "" {
		ext = ".img"
	}
	u := uuid.New().String()
	objectKey := fmt.Sprintf("%s/products/tmp/%s%s", tenantID, u, ext)

	exp := time.Duration(p.cfg.PresignExpiresSec) * time.Second
	uploadURL, publicURL, err := p.storage.PresignPut(context.Background(), p.cfg.R2Bucket, objectKey, contentType, exp)
	if err != nil {
		return nil, domain.InternalServerError{Message: "Error generando URL de subida"}
	}
	return &dto.ProductPhotoInfoDTO{
		UploadURL: uploadURL,
		ObjectKey: objectKey,
		PublicURL: publicURL,
		ExpiresIn: p.cfg.PresignExpiresSec,
	}, nil
}
