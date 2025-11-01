package tenant

import (
	"fmt"

	"gorm.io/gorm"
)

// TenantDB contiene la lógica para manejar schemas por tenant
type TenantDB struct {
	db *gorm.DB
}

// NewTenantDB crea una nueva instancia de TenantDB
func NewTenantDB(db *gorm.DB) *TenantDB {
	return &TenantDB{db: db}
}

// GetTenantSchema retorna el nombre del schema para un tenant específico
func GetTenantSchema(tenantID string) string {
	// Validar que el tenant ID sea válido (alfanumérico y guiones)
	// Para evitar inyección de SQL
	if tenantID == "" {
		return "public"
	}
	return fmt.Sprintf("tenant_%s", tenantID)
}

// SetSchema cambia el schema activo para la conexión actual
func (tdb *TenantDB) SetSchema(tenantID string) *gorm.DB {
	schema := GetTenantSchema(tenantID)
	// Usar el comando SET search_path para cambiar el schema
	return tdb.db.Session(&gorm.Session{}).Exec(fmt.Sprintf("SET search_path TO %s,public", schema))
}

// CreateTenantSchema crea un nuevo schema para un tenant
func (tdb *TenantDB) CreateTenantSchema(tenantID string) error {
	schema := GetTenantSchema(tenantID)

	// Crear el schema
	if err := tdb.db.Exec(fmt.Sprintf("CREATE SCHEMA IF NOT EXISTS %s", schema)).Error; err != nil {
		return fmt.Errorf("error al crear schema para tenant %s: %w", tenantID, err)
	}

	return nil
}

// DeleteTenantSchema elimina el schema de un tenant
func (tdb *TenantDB) DeleteTenantSchema(tenantID string) error {
	schema := GetTenantSchema(tenantID)

	// Eliminar el schema y sus contenidos
	if err := tdb.db.Exec(fmt.Sprintf("DROP SCHEMA IF EXISTS %s CASCADE", schema)).Error; err != nil {
		return fmt.Errorf("error al eliminar schema para tenant %s: %w", tenantID, err)
	}

	return nil
}

// WithTenant retorna una sesión de GORM configurada para trabajar con el schema del tenant
func (tdb *TenantDB) WithTenant(tenantID string) *gorm.DB {
	schema := GetTenantSchema(tenantID)
	return tdb.db.Scopes(func(db *gorm.DB) *gorm.DB {
		// Usar la cláusula de búsqueda de schema en PostgreSQL
		return db.Session(&gorm.Session{NewDB: false}).
			Exec(fmt.Sprintf("SET search_path TO %s,public", schema))
	})
}

// ExecuteInSchema ejecuta una función dentro del contexto de un schema específico
func (tdb *TenantDB) ExecuteInSchema(tenantID string, fn func(*gorm.DB) error) error {
	tx := tdb.db.Begin()
	defer func() {
		if r := recover(); r != nil {
			tx.Rollback()
			panic(r)
		}
	}()

	// Cambiar al schema del tenant
	schema := GetTenantSchema(tenantID)
	if err := tx.Exec(fmt.Sprintf("SET search_path TO %s,public", schema)).Error; err != nil {
		tx.Rollback()
		return err
	}

	// Ejecutar la función
	if err := fn(tx); err != nil {
		tx.Rollback()
		return err
	}

	return tx.Commit().Error
}
