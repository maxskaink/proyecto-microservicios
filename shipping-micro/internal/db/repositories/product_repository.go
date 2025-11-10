package repositories

import (
	"github.com/maxskaink/proyecto-microservicios/shipping-micro/internal/db"
	"github.com/maxskaink/proyecto-microservicios/shipping-micro/internal/db/mappers"
	"github.com/maxskaink/proyecto-microservicios/shipping-micro/internal/db/models"
	"github.com/maxskaink/proyecto-microservicios/shipping-micro/internal/db/tenant"
	"github.com/maxskaink/proyecto-microservicios/shipping-micro/internal/dto"
	"gorm.io/gorm"
)

type ProductRepository struct {
	db       *gorm.DB
	tenantDB *tenant.TenantDB
}

func NewProductRepository(db *gorm.DB, tenantDB *tenant.TenantDB) IProductRepository {
	return &ProductRepository{db: db, tenantDB: tenantDB}
}

func (r *ProductRepository) Create(product *dto.ProductDTO, tenantID string) error {
	productDB := &models.ProductDB{
		ID:          product.ID,
		Name:        product.Name,
		Description: product.Description,
		Price:       product.Price,
		Stock:       product.Stock,
	}

	err := r.tenantDB.ExecuteInSchema(tenantID, func(tx *gorm.DB) error {
		return tx.Create(productDB).Error
	})

	return db.ParseDBError(err)
}

func (r *ProductRepository) GetByID(id string, tenantID string) (*dto.ProductDTO, error) {
	var product models.ProductDB
	err := r.tenantDB.ExecuteInSchema(tenantID, func(tx *gorm.DB) error {
		return tx.Where("id = ?", id).First(&product).Error
	})
	if err != nil {
		return nil, db.ParseDBError(err)
	}
	return mappers.ProductDBToDTO(&product), nil
}

func (r *ProductRepository) Update(product *dto.ProductDTO, tenantID string) error {
	err := r.tenantDB.ExecuteInSchema(tenantID, func(tx *gorm.DB) error {
		return tx.Model(&models.ProductDB{}).
			Where("id = ?", product.ID).
			Updates(map[string]interface{}{
				"name":        product.Name,
				"description": product.Description,
				"price":       product.Price,
				"stock":       product.Stock,
			}).Error
	})

	return db.ParseDBError(err)
}
