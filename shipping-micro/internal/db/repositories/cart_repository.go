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
	var itemID string

	// Crear/actualizar dentro de la transacción y conservar el ID
	if err := r.tenantDB.ExecuteInSchema(tenantID, func(tx *gorm.DB) error {
		var dbItem models.CartItemDB

		// ¿Ya existe el item?
		if err := tx.Where("user_id = ? AND product_id = ?", userID, productID).First(&dbItem).Error; err == nil {
			dbItem.Quantity += quantity
			if err := tx.Save(&dbItem).Error; err != nil {
				return err
			}
			itemID = dbItem.ID
			return nil
		} else if err != gorm.ErrRecordNotFound {
			return err
		}

		// Validar que el producto exista (evita huérfanos)
		var product models.ProductDB
		if err := tx.Where("id = ?", productID).First(&product).Error; err != nil {
			return err
		}

		// Crear nuevo item
		newItem := models.CartItemDB{
			UserID:    userID,
			ProductID: productID,
			Quantity:  quantity,
		}
		if err := tx.Create(&newItem).Error; err != nil {
			return err
		}
		itemID = newItem.ID
		return nil
	}); err != nil {
		return nil, err
	}

	// Fuera de la transacción: recargar con Preload("Product") y devolver
	var full models.CartItemDB
	if err := r.tenantDB.ExecuteInSchema(tenantID, func(tx *gorm.DB) error {
		return tx.Preload("Product").First(&full, "id = ?", itemID).Error
	}); err != nil {
		return nil, err
	}

	return mappers.CartItemDBToDTO(&full), nil
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
