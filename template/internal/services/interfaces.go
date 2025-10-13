package services

import "github.com/maxskaink/proyecto-microservicios/users-micro/internal/domain"

// UserService define la lógica de negocio para usuarios.
type UserService interface {
	CreateUser(u domain.User) (domain.User, error)
	GetUserByID(id string) (domain.User, error)
	UpdateUser(u domain.User) (domain.User, error)
	DeleteUser(id string) error
}

// ProfileService define la lógica de negocio para perfiles de usuario.
type ProfileService interface {
	CreateProfile(p domain.Profile) (domain.Profile, error)
	GetProfileByUserID(userID string) (domain.Profile, error)
	UpdateProfile(p domain.Profile) (domain.Profile, error)
}

// OrderService define la lógica de negocio para pedidos.
type OrderService interface {
	CreateOrder(o domain.Order) (domain.Order, error)
	GetOrderByID(id string) (domain.Order, error)
	ListOrdersByUser(userID string) ([]domain.Order, error)
	UpdateOrderStatus(id string, status string) (domain.Order, error)
}

// AuthService describe operaciones relacionadas con autenticación/autorización.
// La implementación real usará Firebase, pero aquí solo definimos el contrato.
type AuthService interface {
	VerifyToken(token string) (userID string, err error)
}
