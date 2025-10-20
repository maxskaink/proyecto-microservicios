package repositories

import (
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
