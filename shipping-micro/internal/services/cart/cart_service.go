package cart

import (
	"fmt"

	"github.com/maxskaink/proyecto-microservicios/shipping-micro/internal/db/repositories"
	"github.com/maxskaink/proyecto-microservicios/shipping-micro/internal/dto"
)

// Service define la lógica de negocio del carrito
type Service struct {
	cartRepo    repositories.ICartRepository
	productRepo repositories.IProductRepository
}

func NewService(cartRepo repositories.ICartRepository, productRepo repositories.IProductRepository) *Service {
	return &Service{cartRepo: cartRepo, productRepo: productRepo}
}

func (s *Service) AddItem(userID, productID string, quantity int, tenantID string) (*dto.CartItemDTO, error) {
	if quantity <= 0 {
		return nil, fmt.Errorf("quantity must be > 0")
	}
	// Validar producto y stock
	prod, err := s.productRepo.GetByID(productID, tenantID)
	if err != nil {
		return nil, err
	}
	if prod.Stock < quantity {
		// Nota: esto no contempla cantidad previa; se puede mejorar con una consulta del item existente
		return nil, fmt.Errorf("not enough stock for product")
	}
	return s.cartRepo.AddItem(userID, productID, quantity, tenantID)
}

func (s *Service) GetUserCart(userID string, tenantID string) ([]dto.CartItemDTO, error) {
	return s.cartRepo.GetUserCart(userID, tenantID)
}

func (s *Service) UpdateItem(itemID string, quantity int, tenantID string) error {
	if quantity < 0 {
		return fmt.Errorf("quantity must be >= 0")
	}
	// Si vamos a aumentar, validar stock del producto
	if quantity > 0 {
		item, err := s.cartRepo.GetItemByID(itemID, tenantID)
		if err != nil {
			return err
		}
		prod, err := s.productRepo.GetByID(item.ProductID, tenantID)
		if err != nil {
			return err
		}
		if prod.Stock < quantity {
			return fmt.Errorf("not enough stock for product")
		}
	}
	return s.cartRepo.UpdateItem(itemID, quantity, tenantID)
}

func (s *Service) DeleteItem(itemID string, tenantID string) error {
	return s.cartRepo.DeleteItem(itemID, tenantID)
}

func (s *Service) ClearCart(userID string, tenantID string) error {
	return s.cartRepo.ClearCart(userID, tenantID)
}
