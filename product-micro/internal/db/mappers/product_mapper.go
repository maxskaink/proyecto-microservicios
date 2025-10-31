package db_mappers

import (
	db_models "github.com/maxskaink/proyecto-microservicios/product-micro/internal/db/models"
	"github.com/maxskaink/proyecto-microservicios/product-micro/internal/dto"
)

func ProductDtoToModel(dto *dto.ProductDTORequest) *db_models.ProductDB {
	return &db_models.ProductDB{
		ProducerID:  dto.ProducerID,
		Category:    dto.Category,
		Price:       dto.Price,
		Description: dto.Description,
		Stock:       dto.Stock,
		Unit:        dto.Unit,
		PhotoUrl:    dto.PhotoUrl,
	}
}

func ProductModelToDto(model *db_models.ProductDB) *dto.ProductDTOResponse {
	return &dto.ProductDTOResponse{
		ID:          model.ID,
		ProducerID:  model.ProducerID,
		Category:    model.Category,
		Price:       model.Price,
		Description: model.Description,
		Stock:       model.Stock,
		Unit:        model.Unit,
		PhotoUrl:    model.PhotoUrl,
		CreatedAt:   model.CreatedAt,
	}
}
