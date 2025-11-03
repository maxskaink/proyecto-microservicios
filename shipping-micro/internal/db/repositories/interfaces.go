package repositories

import (
	"github.com/maxskaink/proyecto-microservicios/shipping-micro/internal/dto"
)

type IProductRepository interface {
	Create(product *dto.ProductDTO, tenantID string) error
	GetByID(id string, tenantID string) (*dto.ProductDTO, error)
	Update(product *dto.ProductDTO, tenantID string) error
}

type IUserRepository interface {
	Create(user *dto.UserDTO, tenantID string) error
	GetByID(id string, tenantID string) (*dto.UserDTO, error)
	GetByUUID(uid string, tenantID string) (*dto.UserDTO, error)
}

type ICartRepository interface {
	AddItem(userID, productID string, quantity int, tenantID string) (*dto.CartItemDTO, error)
	GetUserCart(userID string, tenantID string) ([]dto.CartItemDTO, error)
	UpdateItem(itemID string, quantity int, tenantID string) error
	DeleteItem(itemID string, tenantID string) error
	ClearCart(userID string, tenantID string) error
	GetItemByID(itemID string, tenantID string) (*dto.CartItemDTO, error)
}

type IOrderRepository interface {
	Create(order *dto.OrderDTO, items []dto.OrderItemDTO, tenantID string) (*dto.OrderDTO, error)
	GetByID(id string, tenantID string) (*dto.OrderDTO, error)
	GetByUserID(userID string, tenantID string) ([]dto.OrderDTO, error)
	UpdateStatus(id string, status string, tenantID string) error
}

type IShippingRepository interface {
	Create(shipping *dto.ShippingDTO, tenantID string) (*dto.ShippingDTO, error)
	GetByID(id string, tenantID string) (*dto.ShippingDTO, error)
	GetByOrderID(orderID string, tenantID string) (*dto.ShippingDTO, error)
	UpdateStatus(id string, status string, tenantID string) error
	List(tenantID string) ([]dto.ShippingDTO, error)
}
