package repositories

import (
	"github.com/maxskaink/proyecto-microservicios/product-micro/internal/db"
	db_mappers "github.com/maxskaink/proyecto-microservicios/product-micro/internal/db/mappers"
	db_models "github.com/maxskaink/proyecto-microservicios/product-micro/internal/db/models"
	"github.com/maxskaink/proyecto-microservicios/product-micro/internal/dto"
	"gorm.io/gorm"
)

// userRepository es la implementación de IUserRepository.
type UserRepository struct {
	db *gorm.DB
}

// NewUserRepository crea una nueva instancia de userRepository.
func NewUserRepository(db *gorm.DB) IUserRepository {
	return &UserRepository{db: db}
}

// CreateUser implements IUserRepository.
func (r *UserRepository) CreateUser(userReq dto.UserRequest) (*dto.UserResponse, error) {
	// Mapear DTO a modelo
	userModel := db_mappers.UserDtoToModel(&userReq)

	// Crear el usuario en la base de datos
	if err := r.db.Create(&userModel).Error; err != nil {
		return nil, db.ParseDBError(err)
	}

	// Mapear el modelo creado a DTO de respuesta
	return db_mappers.UserModelToDto(userModel), nil
}

// FindById implements IUserRepository.
func (r *UserRepository) FindById(id string) (*dto.UserResponse, error) {
	var user db_models.UserDB

	if err := r.db.Preload("Profile").First(&user, "id = ?", id).Error; err != nil {
		return nil, db.ParseDBError(err)
	}

	return db_mappers.UserModelToDto(&user), nil
}

// FindByUID implements IUserRepository.
func (r *UserRepository) FindByUID(uid string) (*dto.UserResponse, error) {
	var user db_models.UserDB

	if err := r.db.First(&user, "firebase_uid = ?", uid).Error; err != nil {
		return nil, db.ParseDBError(err)
	}

	return db_mappers.UserModelToDto(&user), nil
}

// DeleteUser implements IUserRepository.
func (r *UserRepository) DeleteUser(id string) error {
	if err := r.db.Delete(&db_models.UserDB{}, "id = ?", id).Error; err != nil {
		return db.ParseDBError(err)
	}
	return nil
}

// UpdateUser implements IUserRepository.
func (r *UserRepository) UpdateUser(userReq dto.UserRequest) (*dto.UserResponse, error) {
	// Mapear DTO a modelo
	userModel := db_mappers.UserDtoToModel(&userReq)
	// Actualizar el usuario en la base de datos
	if err := r.db.Save(&userModel).Error; err != nil {
		return nil, db.ParseDBError(err)
	}

	// Mapear el modelo actualizado a DTO de respuesta
	return db_mappers.UserModelToDto(userModel), nil
}
