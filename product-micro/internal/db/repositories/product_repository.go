package repositories

import (
	"github.com/maxskaink/proyecto-microservicios/product-micro/internal/db"
	db_mappers "github.com/maxskaink/proyecto-microservicios/product-micro/internal/db/mappers"
	db_models "github.com/maxskaink/proyecto-microservicios/product-micro/internal/db/models"
	"github.com/maxskaink/proyecto-microservicios/product-micro/internal/db/tenant"
	"github.com/maxskaink/proyecto-microservicios/product-micro/internal/dto"
	"gorm.io/gorm"
)

// ProductRepository implementation of IProductRepository.
type ProductRepository struct {
	db       *gorm.DB
	tenantDB *tenant.TenantDB
}

// NewProductRepository create a new instance of ProductRepository.
func NewProductRepository(db *gorm.DB, tenantDB *tenant.TenantDB) IProductRepository {
	return &ProductRepository{
		db:       db,
		tenantDB: tenantDB,
	}
}

// CreateProduct implements IProductRepository.
func (p *ProductRepository) CreateProduct(product *dto.ProductDTORequest, tenantId string) (*dto.ProductDTOResponse, error) {
	productModel := db_mappers.ProductDtoToModel(product)
	var response *dto.ProductDTOResponse

	err := p.tenantDB.ExecuteInSchema(tenantId, func(tx *gorm.DB) error {
		if err := tx.Create(&productModel).Error; err != nil {
			return db.ParseDBError(err)
		}
		response = db_mappers.ProductModelToDto(productModel)
		return nil
	})

	if err != nil {
		return nil, db.ParseDBError(err)
	}

	return response, nil
}

// GetByIdProduct implements IProductRepository.
func (p *ProductRepository) GetByIdProduct(id string, tenantId string) (*dto.ProductDTOResponse, error) {
	var product db_models.ProductDB
	var response *dto.ProductDTOResponse

	err := p.tenantDB.ExecuteInSchema(tenantId, func(tx *gorm.DB) error {
		if err := tx.First(&product, "id = ?", id).Error; err != nil {
			if err == gorm.ErrRecordNotFound {
				return nil
			}
			return db.ParseDBError(err)
		}
		response = db_mappers.ProductModelToDto(&product)
		return nil
	})

	if err != nil {
		return nil, db.ParseDBError(err)
	}

	return response, nil
}

// ListProducts implements IProductRepository con paginación.
func (p *ProductRepository) ListProducts(page int, pageSize int, tenantId string) (*[]dto.ProductDTOResponse, error) {
	if page < 1 {
		page = 1
	}
	if pageSize < 1 {
		pageSize = 10
	}

	var products []db_models.ProductDB
	var responses *[]dto.ProductDTOResponse
	offset := (page - 1) * pageSize

	err := p.tenantDB.ExecuteInSchema(tenantId, func(tx *gorm.DB) error {
		if err := tx.Offset(offset).Limit(pageSize).Find(&products).Error; err != nil {
			return db.ParseDBError(err)
		}

		result := make([]dto.ProductDTOResponse, len(products))
		for i, product := range products {
			result[i] = *db_mappers.ProductModelToDto(&product)
		}
		responses = &result

		return nil
	})

	if err != nil {
		return nil, db.ParseDBError(err)
	}

	return responses, nil
}

// UpdateProduct implements IProductRepository.
func (p *ProductRepository) UpdateProduct(id string, product *dto.ProductDTORequest, tenantId string) (*dto.ProductDTOResponse, error) {
	productModel := db_mappers.ProductDtoToModel(product)
	productModel.ID = id
	var response *dto.ProductDTOResponse

	err := p.tenantDB.ExecuteInSchema(tenantId, func(tx *gorm.DB) error {
		if err := tx.Model(&productModel).Updates(productModel).Error; err != nil {
			return db.ParseDBError(err)
		}
		response = db_mappers.ProductModelToDto(productModel)
		return nil
	})

	if err != nil {
		return nil, db.ParseDBError(err)
	}

	return response, nil
}

// DeleteProduct implements IProductRepository.
func (p *ProductRepository) DeleteProduct(id string, tenantId string) error {
	return p.tenantDB.ExecuteInSchema(tenantId, func(tx *gorm.DB) error {
		if err := tx.Delete(&db_models.ProductDB{}, "id = ?", id).Error; err != nil {
			return db.ParseDBError(err)
		}
		return nil
	})
}
