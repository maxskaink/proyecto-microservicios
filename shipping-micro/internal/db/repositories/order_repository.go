package repositories

import (
	"github.com/maxskaink/proyecto-microservicios/shipping-micro/internal/db/mappers"
	"github.com/maxskaink/proyecto-microservicios/shipping-micro/internal/db/models"
	"github.com/maxskaink/proyecto-microservicios/shipping-micro/internal/db/tenant"
	"github.com/maxskaink/proyecto-microservicios/shipping-micro/internal/dto"
	"gorm.io/gorm"
)

type OrderRepository struct {
	db       *gorm.DB
	tenantDB *tenant.TenantDB
}

func NewOrderRepository(db *gorm.DB, tenantDB *tenant.TenantDB) IOrderRepository {
	return &OrderRepository{db: db, tenantDB: tenantDB}
}

func (r *OrderRepository) Create(order *dto.OrderDTO, items []dto.OrderItemDTO, tenantID string) (*dto.OrderDTO, error) {
	var created models.OrderDB
	err := r.tenantDB.ExecuteInSchema(tenantID, func(tx *gorm.DB) error {
		created = models.OrderDB{
			ID:         order.ID,
			UserID:     order.UserID,
			TotalPrice: order.TotalPrice,
			Status:     order.Status,
		}
		if err := tx.Create(&created).Error; err != nil {
			return err
		}

		// Crear items
		for _, it := range items {
			oi := models.OrderItemDB{
				OrderID:   created.ID,
				ProductID: it.ProductID,
				Quantity:  it.Quantity,
				Price:     it.Price,
			}
			if err := tx.Create(&oi).Error; err != nil {
				return err
			}
			created.Items = append(created.Items, oi)
		}

		return nil
	})
	if err != nil {
		return nil, err
	}
	return mappers.OrderDBToDTO(&created), nil
}

func (r *OrderRepository) GetByID(id string, tenantID string) (*dto.OrderDTO, error) {
	var order models.OrderDB
	err := r.tenantDB.ExecuteInSchema(tenantID, func(tx *gorm.DB) error {
		return tx.Preload("Items").Where("id = ?", id).First(&order).Error
	})
	if err != nil {
		return nil, err
	}
	return mappers.OrderDBToDTO(&order), nil
}

func (r *OrderRepository) GetByUserID(userID string, tenantID string) ([]dto.OrderDTO, error) {
	var orders []models.OrderDB
	err := r.tenantDB.ExecuteInSchema(tenantID, func(tx *gorm.DB) error {
		return tx.Preload("Items").Where("user_id = ?", userID).Order("created_at DESC").Find(&orders).Error
	})
	if err != nil {
		return nil, err
	}
	res := make([]dto.OrderDTO, len(orders))
	for i := range orders {
		res[i] = *mappers.OrderDBToDTO(&orders[i])
	}
	return res, nil
}

func (r *OrderRepository) UpdateStatus(id string, status string, tenantID string) error {
	return r.tenantDB.ExecuteInSchema(tenantID, func(tx *gorm.DB) error {
		return tx.Model(&models.OrderDB{}).Where("id = ?", id).Update("status", status).Error
	})
}
