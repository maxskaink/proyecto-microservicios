package repositories

import (
	"github.com/maxskaink/proyecto-microservicios/tenant-micro/internal/dto"
	"github.com/maxskaink/proyecto-microservicios/tenant-micro/internal/models"
	"gorm.io/gorm"
)

type TenantRepository struct {
	db *gorm.DB
}

func NewTenantRepository(db *gorm.DB) *TenantRepository {
	return &TenantRepository{db: db}
}

func (r *TenantRepository) Create(req dto.CreateTenantRequest) (*dto.TenantResponse, error) {
	tenant := models.Tenant{
		TenantID:    req.TenantID,
		TenantName:  req.TenantName,
		Description: req.Description,
		Location:    req.Location,
	}

	if err := r.db.Create(&tenant).Error; err != nil {
		return nil, err
	}

	return &dto.TenantResponse{
		TenantID:    tenant.TenantID,
		TenantName:  tenant.TenantName,
		Description: tenant.Description,
		Location:    tenant.Location,
		CreatedAt:   tenant.CreatedAt.String(),
	}, nil
}

func (r *TenantRepository) FindByID(tenantID string) (*dto.TenantResponse, error) {
	var tenant models.Tenant
	if err := r.db.Where("tenant_id = ?", tenantID).First(&tenant).Error; err != nil {
		return nil, err
	}

	return &dto.TenantResponse{
		TenantID:    tenant.TenantID,
		TenantName:  tenant.TenantName,
		Description: tenant.Description,
		Location:    tenant.Location,
		CreatedAt:   tenant.CreatedAt.String(),
	}, nil
}

func (r *TenantRepository) GetAll() ([]dto.TenantResponse, error) {
	var tenants []models.Tenant
	if err := r.db.Find(&tenants).Error; err != nil {
		return nil, err
	}

	var responses []dto.TenantResponse
	for _, t := range tenants {
		responses = append(responses, dto.TenantResponse{
			TenantID:    t.TenantID,
			TenantName:  t.TenantName,
			Description: t.Description,
			Location:    t.Location,
			CreatedAt:   t.CreatedAt.String(),
		})
	}

	return responses, nil
}

func (r *TenantRepository) Delete(tenantID string) error {
	return r.db.Where("tenant_id = ?", tenantID).Delete(&models.Tenant{}).Error
}
