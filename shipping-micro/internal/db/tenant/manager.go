package tenant

import (
"fmt"

db_models "github.com/maxskaink/proyecto-microservicios/shipping-micro/internal/db/models"
"gorm.io/gorm"
)

type TenantDB struct {
	db *gorm.DB
}

func NewTenantDB(db *gorm.DB) *TenantDB {
	return &TenantDB{db: db}
}

func GetTenantSchema(tenantID string) string {
	if tenantID == "" {
		return "public"
	}
	return tenantID
}

func (tdb *TenantDB) SetSchema(tenantID string) *gorm.DB {
	schema := GetTenantSchema(tenantID)
	return tdb.db.Session(&gorm.Session{}).Exec(fmt.Sprintf("SET search_path TO %s,public", schema))
}

func (tdb *TenantDB) CreateTenantSchema(tenantID string) error {
	schema := GetTenantSchema(tenantID)

	if err := tdb.db.Exec(fmt.Sprintf("CREATE SCHEMA IF NOT EXISTS %s", schema)).Error; err != nil {
		return fmt.Errorf("error al crear schema para tenant %s: %w", tenantID, err)
	}

	err := tdb.ExecuteInSchema(tenantID, func(tx *gorm.DB) error {
return tx.AutoMigrate(
&db_models.UserDB{},
			&db_models.ProductDB{},
			&db_models.CartItemDB{},
			&db_models.OrderDB{},
			&db_models.OrderItemDB{},
			&db_models.ShippingDB{},
		)
	})

	if err != nil {
		return fmt.Errorf("error al migrar tablas para tenant %s: %w", tenantID, err)
	}

	return nil
}

func (tdb *TenantDB) DeleteTenantSchema(tenantID string) error {
	schema := GetTenantSchema(tenantID)

	if err := tdb.db.Exec(fmt.Sprintf("DROP SCHEMA IF EXISTS %s CASCADE", schema)).Error; err != nil {
		return fmt.Errorf("error al eliminar schema para tenant %s: %w", tenantID, err)
	}

	return nil
}

func (tdb *TenantDB) WithTenant(tenantID string) *gorm.DB {
	schema := GetTenantSchema(tenantID)
	return tdb.db.Scopes(func(db *gorm.DB) *gorm.DB {
return db.Session(&gorm.Session{NewDB: false}).
			Exec(fmt.Sprintf("SET search_path TO %s,public", schema))
	})
}

func (tdb *TenantDB) ExecuteInSchema(tenantID string, fn func(*gorm.DB) error) error {
	tx := tdb.db.Begin()
	defer func() {
		if r := recover(); r != nil {
			tx.Rollback()
			panic(r)
		}
	}()

	schema := GetTenantSchema(tenantID)
	if err := tx.Exec(fmt.Sprintf("SET search_path TO %s,public", schema)).Error; err != nil {
		tx.Rollback()
		return err
	}

	if err := fn(tx); err != nil {
		tx.Rollback()
		return err
	}

	return tx.Commit().Error
}
