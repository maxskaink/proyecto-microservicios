package db_mappers

import (
	db_models "github.com/maxskaink/proyecto-microservicios/users-micro/internal/db/models"
	"github.com/maxskaink/proyecto-microservicios/users-micro/internal/dto"
)

func ProfileModelToDto(model *db_models.ProfileDB) *dto.ProfileResponse {
	return &dto.ProfileResponse{
		ID:          model.ID,
		UserID:      model.UserID,
		Address:     model.Address,
		Phone:       model.Phone,
		AvatarURL:   model.AvatarURL,
		Description: model.Description,
		CreatedAt:   model.CreatedAt,
		UpdatedAt:   model.UpdatedAt,
	}
}

func ProfileDtoToModel(dto *dto.ProfileRequest) *db_models.ProfileDB {
	return &db_models.ProfileDB{
		ID:          dto.UserID,
		UserID:      dto.UserID,
		Address:     dto.Address,
		Description: dto.Description,
		Phone:       dto.Phone,
		AvatarURL:   dto.AvatarURL,
	}
}
