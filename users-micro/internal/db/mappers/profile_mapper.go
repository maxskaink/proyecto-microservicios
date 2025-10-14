package db_mappers

import (
	db_models "github.com/maxskaink/proyecto-microservicios/users-micro/internal/db/models"
	"github.com/maxskaink/proyecto-microservicios/users-micro/internal/dto"
)

func ProfileModelToDto(model *db_models.ProfileDB) *dto.ProfileResponse {
	return &dto.ProfileResponse{
		ID:        model.ID,
		UserID:    model.UserID,
		Address:   model.Address,
		Phone:     model.Phone,
		AvatarURL: model.AvatarURL,
		CreatedAt: model.CreatedAt,
		UpdatedAt: model.UpdatedAt,
	}
}
