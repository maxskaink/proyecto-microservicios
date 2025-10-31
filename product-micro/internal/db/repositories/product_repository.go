package repositories

import (
	db_mappers "github.com/maxskaink/proyecto-microservicios/product-micro/internal/db/mappers"
	db_models "github.com/maxskaink/proyecto-microservicios/product-micro/internal/db/models"
	"github.com/maxskaink/proyecto-microservicios/product-micro/internal/dto"
	"gorm.io/gorm"
)

// ProductRepository implementation of IProductRepository.
type ProductRepository struct {
	db *gorm.DB
}

// NewProductRepository create a new instance of ProductRepository.
func NewProductRepository(db *gorm.DB) IProductRepository {
	return &ProductRepository{db: db}
}

// CreateProduct implements IProductRepository.
func (p *ProductRepository) CreateProduct(product *dto.ProductDTORequest) (*dto.ProductDTOResponse, error) {
	productModel := db_mappers.ProductDtoToModel(product)

	if err := p.db.Create(&productModel).Error; err != nil {
		return nil, err
	}

	response := db_mappers.ProductModelToDto(productModel)
	return response, nil
}

// GetByIdProduct implements IProductRepository.
func (p *ProductRepository) GetByIdProduct(id string) (*dto.ProductDTOResponse, error) {
	var product db_models.ProductDB

	if err := p.db.First(&product, "id = ?", id).Error; err != nil {
		if err == gorm.ErrRecordNotFound {
			return nil, nil
		}
		return nil, err
	}

	response := db_mappers.ProductModelToDto(&product)
	return response, nil
}

// ListProducts implements IProductRepository con paginación.
func (p *ProductRepository) ListProducts(page int, pageSize int) (*[]dto.ProductDTOResponse, error) {

	if page < 1 {
		page = 1
	}
	if pageSize < 1 {
		pageSize = 10
	}

	var products []db_models.ProductDB
	offset := (page - 1) * pageSize

	if err := p.db.Offset(offset).Limit(pageSize).Find(&products).Error; err != nil {
		return nil, err
	}

	responses := make([]dto.ProductDTOResponse, len(products))
	for i, product := range products {
		responses[i] = *db_mappers.ProductModelToDto(&product)
	}

	return &responses, nil
}

// UpdateProduct implements IProductRepository.
func (p *ProductRepository) UpdateProduct(id string, product *dto.ProductDTORequest) (*dto.ProductDTOResponse, error) {
	productModel := db_mappers.ProductDtoToModel(product)
	productModel.ID = id

	if err := p.db.Model(&productModel).Updates(productModel).Error; err != nil {
		return nil, err
	}

	response := db_mappers.ProductModelToDto(productModel)
	return response, nil
}
