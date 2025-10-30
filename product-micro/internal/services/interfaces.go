package services

import (
	"github.com/maxskaink/proyecto-microservicios/product-micro/internal/dto"
)

// IUserService define la lógica de negocio para usuarios.
type IUserService interface {
	GetUserByID(id string) (dto.UserResponse, error)
	GetUserByUUID(id string) (dto.UserResponse, error)
}

type IProductService interface {
}
