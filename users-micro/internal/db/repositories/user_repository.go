package repositories

import (
	"time"

	db_mappers "github.com/maxskaink/proyecto-microservicios/users-micro/internal/db/mappers"
	db_models "github.com/maxskaink/proyecto-microservicios/users-micro/internal/db/models"
	"github.com/maxskaink/proyecto-microservicios/users-micro/internal/dto"
	"gorm.io/gorm"
)

// userRepository es la implementación de UserRepository.
type userRepository struct {
	db *gorm.DB
}

// NewUserRepository crea una nueva instancia de userRepository.
func NewUserRepository(db *gorm.DB) UserRepository {
	return &userRepository{db: db}
}

// Create guarda un nuevo usuario en la base de datos.
func (r *userRepository) Create(u *dto.UserRequest) (*dto.UserResponse, error) {
	var id string
	err := r.db.Transaction(func(tx *gorm.DB) error {
		user := db_mappers.UserDtoToModel(u)

		if err := tx.Create(user).Error; err != nil {
			return err
		}

		profile := &db_models.ProfileDB{
			UserID:    user.ID,
			AvatarURL: "A default avatar",
		}

		if err := tx.Create(profile).Error; err != nil {
			return err
		}
		id = user.ID
		return nil
	})

	if err != nil {
		return nil, err
	}

	userResp, err := r.FindByID(id)
	return userResp, err
}

// FindByID busca un usuario por su ID.
func (r *userRepository) FindByID(id string) (*dto.UserResponse, error) {
	var user db_models.UserDB

	if err := r.db.Preload("Profile").First(&user, "id = ?", id).Error; err != nil {
		return nil, err
	}

	return db_mappers.UserModelToDto(&user), nil
}

func (r *userRepository) FindByUUID(id string) (*dto.UserResponse, error) {
	var user db_models.UserDB

	if err := r.db.Preload("Profile").First(&user, "firebase_uid = ?", id).Error; err != nil {
		return nil, err
	}

	return db_mappers.UserModelToDto(&user), nil
}

// Update actualiza un usuario existente en la base de datos.
func (r *userRepository) Update(id string, u *dto.UserRequest) (*dto.UserResponse, error) {
	var user db_models.UserDB
	if err := r.db.Preload("Profile").First(&user, "id = ?", id).Error; err != nil {
		return nil, err
	}

	user.Name = u.Name
	user.Email = u.Email

	if u.Profile != nil {

		if err := r.db.Model(&user.Profile).Updates(db_models.ProfileDB{
			Address:   u.Profile.Address,
			Phone:     u.Profile.Phone,
			AvatarURL: u.Profile.AvatarURL,
		}).Error; err != nil {
			return nil, err
		}

	}

	user.UpdatedAt = time.Now()

	if err := r.db.Save(&user).Error; err != nil {
		return nil, err
	}

	return db_mappers.UserModelToDto(&user), nil
}

// Delete elimina un usuario por su ID.
func (r *userRepository) Delete(id string) error {
	return r.db.Delete(&db_models.UserDB{}, "id = ?", id).Error
}
