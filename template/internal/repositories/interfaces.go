package repositories

import "github.com/maxskaink/proyecto-microservicios/users-micro/internal/domain"

// UserRepository define operaciones de persistencia para usuarios.
type UserRepository interface {
	Create(u *domain.User) error
	FindByID(id string) (*domain.User, error)
	Update(u *domain.User) error
	Delete(id string) error
}

// ProfileRepository define operaciones de persistencia para perfiles.
type ProfileRepository interface {
	Create(p *domain.Profile) error
	FindByUserID(userID string) (*domain.Profile, error)
	Update(p *domain.Profile) error
}

// OrderRepository define operaciones de persistencia para pedidos.
type OrderRepository interface {
	Create(o *domain.Order) error
	FindByID(id string) (*domain.Order, error)
	ListByUser(userID string) ([]domain.Order, error)
	Update(o *domain.Order) error
}
