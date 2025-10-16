package db_mappers

import (
	db_models "github.com/maxskaink/proyecto-microservicios/users-micro/internal/db/models"
	"github.com/maxskaink/proyecto-microservicios/users-micro/internal/domain"
	"github.com/maxskaink/proyecto-microservicios/users-micro/internal/dto"
)

func UserDtoToModel(dto *dto.UserRequest) *db_models.UserDB {
	return &db_models.UserDB{
		FirebaseUID: dto.FirebaseUID,
		Email:       dto.Email,
		Name:        dto.Name,
		Rol:         string(dto.Rol),
		Profile:     *ProfileDtoToModel(dto.Profile),
	}
}

func UserModelToDto(model *db_models.UserDB) *dto.UserResponse {
	return &dto.UserResponse{
		ID:          model.ID,
		FirebaseUID: model.FirebaseUID,
		Email:       model.Email,
		Name:        model.Name,
		Rol:         domain.UserRole(model.Rol),
		Profile:     *ProfileModelToDto(&model.Profile),
	}
}
