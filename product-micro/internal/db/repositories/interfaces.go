package repositories

import "github.com/maxskaink/proyecto-microservicios/users-micro/internal/dto"

// UserRepository define operaciones de persistencia para usuarios.
type IProductRepository interface {
	// Create
	// FindById
	// Update
	// UpdateROl
	// Delete
}

type IUserRepository interface {
	FindById(id string) (*dto.UserResponse, error)
	FindByUID(uid string) (*dto.UserResponse, error)
}
