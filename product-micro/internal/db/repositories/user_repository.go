package repositories

import (
	"github.com/maxskaink/proyecto-microservicios/product-micro/internal/db"
	db_mappers "github.com/maxskaink/proyecto-microservicios/product-micro/internal/db/mappers"
	db_models "github.com/maxskaink/proyecto-microservicios/product-micro/internal/db/models"
	"github.com/maxskaink/proyecto-microservicios/product-micro/internal/db/tenant"
	"github.com/maxskaink/proyecto-microservicios/product-micro/internal/dto"
	"gorm.io/gorm"
)

// userRepository es la implementación de IUserRepository.
type UserRepository struct {
	db       *gorm.DB
	tenantDB *tenant.TenantDB
}

// NewUserRepository crea una nueva instancia de userRepository.
func NewUserRepository(db *gorm.DB, tenantDB *tenant.TenantDB) IUserRepository {
	return &UserRepository{
		db:       db,
		tenantDB: tenantDB,
	}
}

// CreateUser implements IUserRepository.
func (r *UserRepository) CreateUser(userReq dto.UserRequest, tenantId string) (*dto.UserResponse, error) {
	// Mapear DTO a modelo
	userModel := db_mappers.UserDtoToModel(&userReq)

	// Crear el usuario en el schema del tenant
	var result *dto.UserResponse
	err := r.tenantDB.ExecuteInSchema(tenantId, func(tx *gorm.DB) error {
		if err := tx.Create(&userModel).Error; err != nil {
			return db.ParseDBError(err)
		}
		result = db_mappers.UserModelToDto(userModel)
		return nil
	})

	if err != nil {
		return nil, err
	}

	return result, nil
}

// FindById implements IUserRepository.
func (r *UserRepository) FindById(id string, tenantId string) (*dto.UserResponse, error) {
	var user db_models.UserDB
	var result *dto.UserResponse

	err := r.tenantDB.ExecuteInSchema(tenantId, func(tx *gorm.DB) error {
		if err := tx.Preload("Profile").First(&user, "id = ?", id).Error; err != nil {
			return db.ParseDBError(err)
		}
		result = db_mappers.UserModelToDto(&user)
		return nil
	})

	if err != nil {
		return nil, err
	}

	return result, nil
}

// FindByUID implements IUserRepository.
func (r *UserRepository) FindByUID(uid string, tenantId string) (*dto.UserResponse, error) {
	var user db_models.UserDB
	var result *dto.UserResponse

	err := r.tenantDB.ExecuteInSchema(tenantId, func(tx *gorm.DB) error {
		if err := tx.First(&user, "firebase_uid = ?", uid).Error; err != nil {
			return db.ParseDBError(err)
		}
		result = db_mappers.UserModelToDto(&user)
		return nil
	})

	if err != nil {
		return nil, err
	}

	return result, nil
}

// DeleteUser implements IUserRepository.
func (r *UserRepository) DeleteUser(id string, tenantId string) error {
	return r.tenantDB.ExecuteInSchema(tenantId, func(tx *gorm.DB) error {
		if err := tx.Delete(&db_models.UserDB{}, "id = ?", id).Error; err != nil {
			return db.ParseDBError(err)
		}
		return nil
	})
}

// UpdateUser implements IUserRepository.
func (r *UserRepository) UpdateUser(userReq dto.UserRequest, tenantId string) (*dto.UserResponse, error) {
	// Mapear DTO a modelo
	userModel := db_mappers.UserDtoToModel(&userReq)
	var result *dto.UserResponse

	err := r.tenantDB.ExecuteInSchema(tenantId, func(tx *gorm.DB) error {
		// Actualizar el usuario en el schema del tenant
		if err := tx.Save(&userModel).Error; err != nil {
			return db.ParseDBError(err)
		}
		result = db_mappers.UserModelToDto(userModel)
		return nil
	})

	if err != nil {
		return nil, err
	}

	return result, nil
}
