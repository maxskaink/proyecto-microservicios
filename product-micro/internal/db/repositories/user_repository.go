package repositories

import (
	"github.com/maxskaink/proyecto-microservicios/users-micro/internal/db"
	db_mappers "github.com/maxskaink/proyecto-microservicios/users-micro/internal/db/mappers"
	db_models "github.com/maxskaink/proyecto-microservicios/users-micro/internal/db/models"
	"github.com/maxskaink/proyecto-microservicios/users-micro/internal/dto"
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

	if err := r.db.Preload("Profile").First(&user, "firebase_uid = ?", uid).Error; err != nil {
		return nil, db.ParseDBError(err)
	}

	return db_mappers.UserModelToDto(&user), nil
}
