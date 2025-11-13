package repositories

import (
	"fmt"

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

// Update actualiza los datos del usuario identificado por ID, UUID o Email (en ese orden)
func (r *UserRepository) Update(user *dto.UserDTO, tenantID string) error {
	err := r.tenantDB.ExecuteInSchema(tenantID, func(tx *gorm.DB) error {
		// Determinar criterio de búsqueda
		var where string
		var arg interface{}
		switch {
		case user.ID != "":
			where, arg = "id = ?", user.ID
		case user.FirebaseUID != "":
			where, arg = "uuid = ?", user.FirebaseUID
		case user.Email != "":
			where, arg = "email = ?", user.Email
		default:
			return fmt.Errorf("no identifier provided to update user")
		}

		updates := map[string]interface{}{}
		if user.Email != "" {
			updates["email"] = user.Email
		}
		if user.Name != "" {
			updates["name"] = user.Name
		}
		if user.Rol != "" {
			updates["rol"] = user.Rol
		}
		if user.FirebaseUID != "" {
			updates["uuid"] = user.FirebaseUID
		}
		if len(updates) == 0 {
			// Nada que actualizar
			return nil
		}
		return tx.Model(&models.UserDB{}).Where(where, arg).Updates(updates).Error
	})

	return db.ParseDBError(err)
}
