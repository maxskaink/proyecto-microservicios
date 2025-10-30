package db_mappers

import (
	db_models "github.com/maxskaink/proyecto-microservicios/product-micro/internal/db/models"
	"github.com/maxskaink/proyecto-microservicios/product-micro/internal/dto"
)

func UserDtoToModel(dto *dto.UserRequest) *db_models.UserDB {
	return &db_models.UserDB{
		FirebaseUID: dto.FirebaseUID,
		Email:       dto.Email,
		Name:        dto.Name,
		Rol:         string(dto.Rol),
	}
}

func UserModelToDto(model *db_models.UserDB) *dto.UserResponse {
	return &dto.UserResponse{
		ID:          model.ID,
		FirebaseUID: model.FirebaseUID,
		Email:       model.Email,
		Name:        model.Name,
		Rol:         model.Rol,
	}
}
