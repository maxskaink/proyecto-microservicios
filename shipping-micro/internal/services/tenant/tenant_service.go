package tenant

import (
"github.com/maxskaink/proyecto-microservicios/shipping-micro/internal/db/tenant"
"gorm.io/gorm"
)

type TenantService struct {
	db       *gorm.DB
	tenantDB *tenant.TenantDB
}

func NewTenantService(db *gorm.DB) *TenantService {
	return &TenantService{
		db:       db,
		tenantDB: tenant.NewTenantDB(db),
	}
}

func (ts *TenantService) CreateTenantSchema(tenantID string) error {
	return ts.tenantDB.CreateTenantSchema(tenantID)
}

func (ts *TenantService) DeleteTenantSchema(tenantID string) error {
	return ts.tenantDB.DeleteTenantSchema(tenantID)
}
