package order

import (
	"context"
	"errors"
	"fmt"

	"github.com/google/uuid"
	"github.com/maxskaink/proyecto-microservicios/shipping-micro/internal/db/repositories"
	"github.com/maxskaink/proyecto-microservicios/shipping-micro/internal/dto"
	"github.com/maxskaink/proyecto-microservicios/shipping-micro/internal/messaging"
)

// Service contiene la lógica de órdenes
type Service struct {
	cartRepo     repositories.ICartRepository
	productRepo  repositories.IProductRepository
	orderRepo    repositories.IOrderRepository
	shippingRepo repositories.IShippingRepository
	publisher    messaging.Publisher
}

func NewService(cart repositories.ICartRepository, prod repositories.IProductRepository, ord repositories.IOrderRepository, ship repositories.IShippingRepository, pub messaging.Publisher) *Service {
	return &Service{cartRepo: cart, productRepo: prod, orderRepo: ord, shippingRepo: ship, publisher: pub}
}

// CreateFromCart crea una orden desde el carrito del usuario y limpia el carrito.
func (s *Service) CreateFromCart(userID string, shippingAddress string, tenantID string) (*dto.OrderDTO, *dto.ShippingDTO, error) {
	items, err := s.cartRepo.GetUserCart(userID, tenantID)
	if err != nil {
		return nil, nil, err
	}
	if len(items) == 0 {
		return nil, nil, errors.New("cart is empty")
	}

	// Calcular total y mapear items
	var total float64
	var orderItems []dto.OrderItemDTO
	for _, it := range items {
		prod, err := s.productRepo.GetByID(it.ProductID, tenantID)
		if err != nil {
			return nil, nil, err
		}
		total += float64(it.Quantity) * prod.Price
		orderItems = append(orderItems, dto.OrderItemDTO{
			ProductID: it.ProductID,
			Quantity:  it.Quantity,
			Price:     prod.Price,
		})
	}

	order := &dto.OrderDTO{UserID: userID, TotalPrice: total, Status: "pending"}
	createdOrder, err := s.orderRepo.Create(order, orderItems, tenantID)
	if err != nil {
		return nil, nil, err
	}

	// Crear shipping en estado pending con tracking generado
	shipping := &dto.ShippingDTO{
		OrderID:         createdOrder.ID,
		TrackingNumber:  generateTracking(),
		ShippingAddress: shippingAddress,
		Status:          "pending",
	}
	createdShipping, err := s.shippingRepo.Create(shipping, tenantID)
	if err != nil {
		return nil, nil, err
	}

	// Limpiar carrito
	if err := s.cartRepo.ClearCart(userID, tenantID); err != nil {
		return nil, nil, fmt.Errorf("order created but failed to clear cart: %w", err)
	}

	// Publicar eventos
	if s.publisher != nil {
		ctx := context.Background()
		_ = s.publisher.PublishOrderCreated(ctx, map[string]interface{}{
			"id":          createdOrder.ID,
			"user_id":     createdOrder.UserID,
			"total_price": createdOrder.TotalPrice,
			"status":      createdOrder.Status,
		}, tenantID)
		_ = s.publisher.PublishShippingCreated(ctx, map[string]interface{}{
			"id":              createdShipping.ID,
			"order_id":        createdShipping.OrderID,
			"tracking_number": createdShipping.TrackingNumber,
			"status":          createdShipping.Status,
		}, tenantID)
	}

	return createdOrder, createdShipping, nil
}

func (s *Service) GetByID(orderID string, tenantID string) (*dto.OrderDTO, error) {
	return s.orderRepo.GetByID(orderID, tenantID)
}

func (s *Service) GetByUserID(userID string, tenantID string) ([]dto.OrderDTO, error) {
	return s.orderRepo.GetByUserID(userID, tenantID)
}

func (s *Service) UpdateStatus(orderID string, status string, tenantID string) error {
	if status != "pending" && status != "paid" && status != "cancelled" {
		return fmt.Errorf("invalid status")
	}
	if err := s.orderRepo.UpdateStatus(orderID, status, tenantID); err != nil {
		return err
	}
	// Si la orden pasa a paid, avanzar el shipping a in_transit si existe
	if status == "paid" {
		if sh, err := s.shippingRepo.GetByOrderID(orderID, tenantID); err == nil && sh != nil {
			_ = s.shippingRepo.UpdateStatus(sh.ID, "in_transit", tenantID)
		}
	}
	if s.publisher != nil {
		ctx := context.Background()
		_ = s.publisher.PublishOrderStatusChanged(ctx, map[string]interface{}{
			"id":     orderID,
			"status": status,
		}, tenantID)
	}
	return nil
}

// generateTracking crea un tracking number simple
func generateTracking() string {
	return "TRK-" + uuid.New().String()[0:8]
}
