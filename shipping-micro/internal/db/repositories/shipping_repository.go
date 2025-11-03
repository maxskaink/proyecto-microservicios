package repositories

import (
	"github.com/maxskaink/proyecto-microservicios/shipping-micro/internal/db/mappers"
	"github.com/maxskaink/proyecto-microservicios/shipping-micro/internal/db/models"
	"github.com/maxskaink/proyecto-microservicios/shipping-micro/internal/db/tenant"
	"github.com/maxskaink/proyecto-microservicios/shipping-micro/internal/dto"
	"gorm.io/gorm"
)

type ShippingRepository struct {
	db       *gorm.DB
	tenantDB *tenant.TenantDB
}

func NewShippingRepository(db *gorm.DB, tenantDB *tenant.TenantDB) IShippingRepository {
	return &ShippingRepository{db: db, tenantDB: tenantDB}
}

func (r *ShippingRepository) Create(shipping *dto.ShippingDTO, tenantID string) (*dto.ShippingDTO, error) {
	var created models.ShippingDB
	err := r.tenantDB.ExecuteInSchema(tenantID, func(tx *gorm.DB) error {
		created = models.ShippingDB{
			ID:              shipping.ID,
			OrderID:         shipping.OrderID,
			TrackingNumber:  shipping.TrackingNumber,
			ShippingAddress: shipping.ShippingAddress,
			Status:          shipping.Status,
		}
		return tx.Create(&created).Error
	})
	if err != nil {
		return nil, err
	}
	return mappers.ShippingDBToDTO(&created), nil
}

func (r *ShippingRepository) GetByID(id string, tenantID string) (*dto.ShippingDTO, error) {
	var s models.ShippingDB
	err := r.tenantDB.ExecuteInSchema(tenantID, func(tx *gorm.DB) error {
		return tx.Where("id = ?", id).First(&s).Error
	})
	if err != nil {
		return nil, err
	}
	return mappers.ShippingDBToDTO(&s), nil
}

func (r *ShippingRepository) GetByOrderID(orderID string, tenantID string) (*dto.ShippingDTO, error) {
	var s models.ShippingDB
	err := r.tenantDB.ExecuteInSchema(tenantID, func(tx *gorm.DB) error {
		return tx.Where("order_id = ?", orderID).First(&s).Error
	})
	if err != nil {
		return nil, err
	}
	return mappers.ShippingDBToDTO(&s), nil
}

func (r *ShippingRepository) UpdateStatus(id string, status string, tenantID string) error {
	return r.tenantDB.ExecuteInSchema(tenantID, func(tx *gorm.DB) error {
		return tx.Model(&models.ShippingDB{}).Where("id = ?", id).Update("status", status).Error
	})
}

func (r *ShippingRepository) List(tenantID string) ([]dto.ShippingDTO, error) {
	var sh []models.ShippingDB
	err := r.tenantDB.ExecuteInSchema(tenantID, func(tx *gorm.DB) error {
		return tx.Order("created_at DESC").Find(&sh).Error
	})
	if err != nil {
		return nil, err
	}
	res := make([]dto.ShippingDTO, len(sh))
	for i := range sh {
		res[i] = *mappers.ShippingDBToDTO(&sh[i])
	}
	return res, nil
}
