package repositories

import (
	"fmt"
	"time"

	"github.com/maxskaink/proyecto-microservicios/users-micro/internal/db"
	db_mappers "github.com/maxskaink/proyecto-microservicios/users-micro/internal/db/mappers"
	db_models "github.com/maxskaink/proyecto-microservicios/users-micro/internal/db/models"
	"github.com/maxskaink/proyecto-microservicios/users-micro/internal/db/tenant"
	"github.com/maxskaink/proyecto-microservicios/users-micro/internal/domain"
	"github.com/maxskaink/proyecto-microservicios/users-micro/internal/dto"
	"gorm.io/gorm"
)

// userRepository es la implementación de UserRepository.
type userRepository struct {
	db     *gorm.DB
	tenant tenant.TenantDB
}

// NewUserRepository crea una nueva instancia de userRepository.
func NewUserRepository(db *gorm.DB, tenant tenant.TenantDB) UserRepository {
	return &userRepository{db: db, tenant: tenant}
}

// Create guarda un nuevo usuario en la base de datos.
func (r *userRepository) Create(u *dto.UserRequest, tenantID string) (*dto.UserResponse, error) {
	var id string

	fmt.Printf("Creating user in tenant schema: %s\n", tenantID)

	err := r.tenant.ExecuteInSchema(tenantID, func(tx *gorm.DB) error {
		return tx.Transaction(func(txn *gorm.DB) error {
			user := db_mappers.UserDtoToModel(u)

			var count int64
			if err := txn.Model(&db_models.UserDB{}).Count(&count).Error; err != nil {
				return err
			}

			if count == 0 {
				user.Rol = string(domain.UserRoleAdmin)
			}

			// Evitar que GORM intente crear la relación Profile automáticamente
			if err := txn.Omit("Profile").Create(user).Error; err != nil {
				return err
			}

			profile := &db_models.ProfileDB{
				UserID:    user.ID,
				AvatarURL: "A default avatar",
			}

			if err := txn.Create(profile).Error; err != nil {
				return err
			}
			id = user.ID
			return nil
		})
	})

	if err != nil {
		return nil, db.ParseDBError(err)
	}

	userResp, err := r.FindByID(id, tenantID)
	return userResp, err
}

// FindByID busca un usuario por su ID.
func (r *userRepository) FindByID(id string, tenantID string) (*dto.UserResponse, error) {
	var user db_models.UserDB

	err := r.tenant.ExecuteInSchema(tenantID, func(tx *gorm.DB) error {
		return tx.Preload("Profile").First(&user, "id = ?", id).Error
	})

	if err != nil {
		return nil, db.ParseDBError(err)
	}

	return db_mappers.UserModelToDto(&user), nil
}

func (r *userRepository) FindByUUID(id string, tenantID string) (*dto.UserResponse, error) {
	var user db_models.UserDB

	err := r.tenant.ExecuteInSchema(tenantID, func(tx *gorm.DB) error {
		return tx.Preload("Profile").First(&user, "firebase_uid = ?", id).Error
	})

	if err != nil {
		return nil, db.ParseDBError(err)
	}

	return db_mappers.UserModelToDto(&user), nil
}

// Update actualiza un usuario existente en la base de datos.
func (r *userRepository) Update(id string, u *dto.UserRequest, tenantID string) (*dto.UserResponse, error) {
	var user db_models.UserDB

	err := r.tenant.ExecuteInSchema(tenantID, func(tx *gorm.DB) error {
		if err := tx.Preload("Profile").First(&user, "id = ?", id).Error; err != nil {
			return err
		}

		user.Name = u.Name
		user.Email = u.Email

		if u.Profile != nil {

			if err := tx.Model(&user.Profile).Updates(db_models.ProfileDB{
				Address:     u.Profile.Address,
				Phone:       u.Profile.Phone,
				Description: u.Profile.Description,
				AvatarURL:   u.Profile.AvatarURL,
			}).Error; err != nil {
				return db.ParseDBError(err)
			}

		}

		user.UpdatedAt = time.Now()

		if err := tx.Save(&user).Error; err != nil {
			return db.ParseDBError(err)
		}

		return nil
	})

	if err != nil {
		return nil, err
	}

	return db_mappers.UserModelToDto(&user), nil
}

// Update the rol of a user
func (r *userRepository) UpdateRol(id string, rol string, tenantID string) (*dto.UserResponse, error) {

	err := r.tenant.ExecuteInSchema(tenantID, func(tx *gorm.DB) error {
		return tx.Model(&db_models.UserDB{}).Where("id = ?", id).Update("rol", rol).Error
	})

	if err != nil {
		return nil, db.ParseDBError(err)
	}

	return r.FindByID(id, tenantID)
}

// Delete elimina un usuario por su ID.
func (r *userRepository) Delete(id string, tenantID string) error {
	return db.ParseDBError(r.tenant.ExecuteInSchema(tenantID, func(tx *gorm.DB) error {
		return tx.Delete(&db_models.UserDB{}, "id = ?", id).Error
	}))
}
