package repositories

import (
	"github.com/maxskaink/proyecto-microservicios/shipping-micro/internal/db/mappers"
	"github.com/maxskaink/proyecto-microservicios/shipping-micro/internal/db/models"
	"github.com/maxskaink/proyecto-microservicios/shipping-micro/internal/db/tenant"
	"github.com/maxskaink/proyecto-microservicios/shipping-micro/internal/dto"
	"gorm.io/gorm"
)

type CartRepository struct {
	db       *gorm.DB
	tenantDB *tenant.TenantDB
}

func NewCartRepository(db *gorm.DB, tenantDB *tenant.TenantDB) ICartRepository {
	return &CartRepository{db: db, tenantDB: tenantDB}
}

func (r *CartRepository) AddItem(userID, productID string, quantity int, tenantID string) (*dto.CartItemDTO, error) {
	var cartItem models.CartItemDB

	err := r.tenantDB.ExecuteInSchema(tenantID, func(tx *gorm.DB) error {
		// Verificar si el item ya existe en el carrito
		err := tx.Where("user_id = ? AND product_id = ?", userID, productID).First(&cartItem).Error
		if err == nil {
			// El item ya existe, actualizar cantidad
			cartItem.Quantity += quantity
			return tx.Save(&cartItem).Error
		}

		// El item no existe, crear uno nuevo
		if err != gorm.ErrRecordNotFound {
			return err
		}

		cartItem = models.CartItemDB{
			UserID:    userID,
			ProductID: productID,
			Quantity:  quantity,
		}
		return tx.Create(&cartItem).Error
	})

	if err != nil {
		return nil, err
	}

	return mappers.CartItemDBToDTO(&cartItem), nil
}

func (r *CartRepository) GetUserCart(userID string, tenantID string) ([]dto.CartItemDTO, error) {
	var items []models.CartItemDB

	err := r.tenantDB.ExecuteInSchema(tenantID, func(tx *gorm.DB) error {
		return tx.Preload("Product").Where("user_id = ?", userID).Find(&items).Error
	})

	if err != nil {
		return nil, err
	}

	result := make([]dto.CartItemDTO, len(items))
	for i, item := range items {
		result[i] = *mappers.CartItemDBToDTO(&item)
	}

	return result, nil
}

func (r *CartRepository) UpdateItem(itemID string, quantity int, tenantID string) error {
	return r.tenantDB.ExecuteInSchema(tenantID, func(tx *gorm.DB) error {
		if quantity == 0 {
			// Si la cantidad es 0, eliminar el item
			return tx.Where("id = ?", itemID).Delete(&models.CartItemDB{}).Error
		}
		return tx.Model(&models.CartItemDB{}).Where("id = ?", itemID).Update("quantity", quantity).Error
	})
}

func (r *CartRepository) DeleteItem(itemID string, tenantID string) error {
	return r.tenantDB.ExecuteInSchema(tenantID, func(tx *gorm.DB) error {
		return tx.Where("id = ?", itemID).Delete(&models.CartItemDB{}).Error
	})
}

func (r *CartRepository) ClearCart(userID string, tenantID string) error {
	return r.tenantDB.ExecuteInSchema(tenantID, func(tx *gorm.DB) error {
		return tx.Where("user_id = ?", userID).Delete(&models.CartItemDB{}).Error
	})
}

func (r *CartRepository) GetItemByID(itemID string, tenantID string) (*dto.CartItemDTO, error) {
	var item models.CartItemDB

	err := r.tenantDB.ExecuteInSchema(tenantID, func(tx *gorm.DB) error {
		return tx.Where("id = ?", itemID).First(&item).Error
	})

	if err != nil {
		return nil, err
	}

	return mappers.CartItemDBToDTO(&item), nil
}
