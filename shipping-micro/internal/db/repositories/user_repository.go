package repositories

import (
	"github.com/maxskaink/proyecto-microservicios/shipping-micro/internal/db"
	"github.com/maxskaink/proyecto-microservicios/shipping-micro/internal/db/mappers"
	"github.com/maxskaink/proyecto-microservicios/shipping-micro/internal/db/models"
	"github.com/maxskaink/proyecto-microservicios/shipping-micro/internal/db/tenant"
	"github.com/maxskaink/proyecto-microservicios/shipping-micro/internal/dto"
	"gorm.io/gorm"
)

type UserRepository struct {
	db       *gorm.DB
	tenantDB *tenant.TenantDB
}

func NewUserRepository(db *gorm.DB, tenantDB *tenant.TenantDB) IUserRepository {
	return &UserRepository{db: db, tenantDB: tenantDB}
}

func (r *UserRepository) Create(user *dto.UserDTO, tenantID string) error {
	userDB := &models.UserDB{
		ID:    user.ID,
		UUID:  user.FirebaseUID,
		Email: user.Email,
		Name:  user.Name,
	}

	err := r.tenantDB.ExecuteInSchema(tenantID, func(tx *gorm.DB) error {
		return tx.Create(userDB).Error
	})

	return db.ParseDBError(err)
}

func (r *UserRepository) GetByID(id string, tenantID string) (*dto.UserDTO, error) {
	var user models.UserDB
	err := r.tenantDB.ExecuteInSchema(tenantID, func(tx *gorm.DB) error {
		return tx.Where("id = ?", id).First(&user).Error
	})
	if err != nil {
		return nil, db.ParseDBError(err)
	}
	return mappers.UserDBToDTO(&user), nil
}

// GetByUUID implements IUserRepository.
func (r *UserRepository) GetByUUID(uid string, tenantID string) (*dto.UserDTO, error) {
	var user models.UserDB
	err := r.tenantDB.ExecuteInSchema(tenantID, func(tx *gorm.DB) error {
		return tx.Where("uuid = ?", uid).First(&user).Error
	})
	if err != nil {
		return nil, db.ParseDBError(err)
	}
	return mappers.UserDBToDTO(&user), nil
}
