package tenant

import (
	"context"
	"fmt"

	"github.com/maxskaink/proyecto-microservicios/product-micro/internal/db/tenant"
	"github.com/maxskaink/proyecto-microservicios/product-micro/pkg/logger"
	"gorm.io/gorm"
)

// TenantService maneja la lógica de tenants en el microservicio de productos
type TenantService struct {
	db     *gorm.DB
	tenant *tenant.TenantDB
}

// NewTenantService crea una nueva instancia de TenantService
func NewTenantService(db *gorm.DB) *TenantService {
	return &TenantService{
		db:     db,
		tenant: tenant.NewTenantDB(db),
	}
}

// CreateTenant crea un nuevo tenant
func (ts *TenantService) CreateTenant(ctx context.Context, tenantID string, tenantName string) error {
	logger.Info(fmt.Sprintf("Creando schema para tenant: %s (%s)", tenantID, tenantName))

	// Crear el schema en la base de datos
	if err := ts.tenant.CreateTenantSchema(tenantID); err != nil {
		logger.Error(fmt.Sprintf("Error al crear schema para tenant %s: %v", tenantID, err))
		return err
	}

	logger.Info(fmt.Sprintf("Schema creado exitosamente para tenant: %s", tenantID))

	return nil
}

// DeleteTenant elimina un tenant
func (ts *TenantService) DeleteTenant(ctx context.Context, tenantID string) error {
	logger.Info(fmt.Sprintf("Eliminando schema para tenant: %s", tenantID))

	// Eliminar el schema de la base de datos
	if err := ts.tenant.DeleteTenantSchema(tenantID); err != nil {
		logger.Error(fmt.Sprintf("Error al eliminar schema para tenant %s: %v", tenantID, err))
		return err
	}

	logger.Info(fmt.Sprintf("Schema eliminado exitosamente para tenant: %s", tenantID))

	return nil
}

// GetTenantSchema obtiene el nombre del schema para un tenant
func (ts *TenantService) GetTenantSchema(tenantID string) string {
	return tenant.GetTenantSchema(tenantID)
}
