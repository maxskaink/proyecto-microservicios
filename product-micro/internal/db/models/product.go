package db_models

import (
	"time"

	"github.com/maxskaink/proyecto-microservicios/product-micro/internal/domain"
)

type ProductDB struct {
	ID          string                 `gorm:"primaryKey;type:uuid;default:gen_random_uuid()"`
	ProducerID  string                 `gorm:"type:uuid;not null;index"`
	Producer    *UserDB                `gorm:"foreignKey:ProducerID;references:ID"`
	Category    domain.ProductCategory `gorm:"type:varchar(100);not null" json:"category"`
	Price       int                    `gorm:"not null;check:price > 0" json:"price"`
	Description string                 `gorm:"type:text" json:"description"`
	Stock       int                    `gorm:"not null;default:0;check:stock >= 0" json:"stock"`
	Unit        domain.ProductUnit     `gorm:"type:varchar(50)" json:"unit"`
	PhotoUrl    string                 `gorm:"type:text" json:"photo_url"`
	CreatedAt   time.Time              `gorm:"autoCreateTime" json:"created_at"`
	UpdatedAt   time.Time              `gorm:"autoUpdateTime" json:"updated_at"`
	DeletedAt   *time.Time             `gorm:"index" json:"deleted_at"`
}

// TableName especifica el nombre de la tabla
func (ProductDB) TableName() string {
	return "products"
}
