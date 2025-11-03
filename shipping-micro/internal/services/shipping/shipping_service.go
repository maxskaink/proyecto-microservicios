package shipping

import (
	"fmt"

	"github.com/maxskaink/proyecto-microservicios/shipping-micro/internal/db/repositories"
	"github.com/maxskaink/proyecto-microservicios/shipping-micro/internal/dto"
)

// Service maneja la lógica de envíos
type Service struct {
	shippingRepo repositories.IShippingRepository
}

func NewService(repo repositories.IShippingRepository) *Service {
	return &Service{shippingRepo: repo}
}

func (s *Service) GetByID(id string, tenantID string) (*dto.ShippingDTO, error) {
	return s.shippingRepo.GetByID(id, tenantID)
}

func (s *Service) GetByOrderID(orderID string, tenantID string) (*dto.ShippingDTO, error) {
	return s.shippingRepo.GetByOrderID(orderID, tenantID)
}

func (s *Service) List(tenantID string) ([]dto.ShippingDTO, error) {
	return s.shippingRepo.List(tenantID)
}

func (s *Service) UpdateStatus(id string, status string, tenantID string) error {
	switch status {
	case "pending", "in_transit", "delivered", "cancelled":
	default:
		return fmt.Errorf("invalid status")
	}
	return s.shippingRepo.UpdateStatus(id, status, tenantID)
}
