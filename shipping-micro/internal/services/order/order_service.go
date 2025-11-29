package order

import (
	"context"
	"fmt"

	"github.com/google/uuid"
	"github.com/maxskaink/proyecto-microservicios/shipping-micro/internal/db/repositories"
	"github.com/maxskaink/proyecto-microservicios/shipping-micro/internal/domain"
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
	userRepo     repositories.IUserRepository
}

func NewService(cart repositories.ICartRepository, prod repositories.IProductRepository, ord repositories.IOrderRepository, ship repositories.IShippingRepository, pub messaging.Publisher, user repositories.IUserRepository) *Service {
	return &Service{cartRepo: cart, productRepo: prod, orderRepo: ord, shippingRepo: ship, publisher: pub, userRepo: user}
}

// CreateFromCart crea una orden desde el carrito del usuario y limpia el carrito.
func (s *Service) CreateFromCart(userID string, shippingAddress string, tenantID string) (*dto.OrderDTO, *dto.ShippingDTO, error) {
	items, err := s.cartRepo.GetUserCart(userID, tenantID)
	if err != nil {
		return nil, nil, err
	}
	if len(items) == 0 {
		return nil, nil, domain.ConflictError{Message: "cart is empty"}
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

	order := &dto.OrderDTO{UserID: userID, TotalPrice: total, Status: string(domain.OrderStatusPending)}
	createdOrder, err := s.orderRepo.Create(order, orderItems, tenantID)
	if err != nil {
		return nil, nil, err
	}

	// Crear shipping en estado pending con tracking generado
	shipping := &dto.ShippingDTO{
		OrderID:         createdOrder.ID,
		TrackingNumber:  generateTracking(),
		ShippingAddress: shippingAddress,
		Status:          string(domain.OrderStatusPending),
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

func (s *Service) GetByUserID(userID string, status domain.OrderStatus, tenantID string) ([]dto.OrderDTO, error) {

	_, err := s.userRepo.GetByID(userID, tenantID)
	if err != nil {
		return nil, err
	}
	response_raw, err := s.orderRepo.GetByUserID(userID, tenantID)
	if err != nil {
		return nil, err
	}

	//Filtrar por estado
	var response []dto.OrderDTO
	for _, order := range response_raw {
		if order.Status == string(status) {
			response = append(response, order)
		}
	}

	return response, nil

}

func (s *Service) UpdateStatus(orderID string, status string, requester string, tenantID string) error {
	fmt.Println("Actualizanod la orden con id", orderID)
	if status != "pending" && status != "paid" && status != "cancelled" {
		return fmt.Errorf("invalid status")
	}
	//Validate permisiones
	user, err := s.userRepo.GetByID(requester, tenantID)
	if err != nil {
		return err
	}
	order, err := s.orderRepo.GetByID(orderID, tenantID)
	if err != nil {
		return err
	}
	if user.Rol != "admin" && user.ID != order.UserID {
		return domain.UnauthorizedError{Message: "no permission to update order status"}
	}

	// Actualizar estado de la orden

	if err := s.orderRepo.UpdateStatus(orderID, status, tenantID); err != nil {
		return err
	}
	// Si la orden pasa a paid, avanzar el shipping a in_transit si existe
	if status == "paid" {
		if sh, err := s.shippingRepo.GetByOrderID(orderID, tenantID); err == nil && sh != nil {
			_ = s.shippingRepo.UpdateStatus(sh.ID, "in_transit", tenantID)
		}
	}
	//Si la orden es cancelada, actualizar el shipping a cancelled si existe
	if status == "cancelled" {
		if sh, err := s.shippingRepo.GetByOrderID(orderID, tenantID); err == nil && sh != nil {
			_ = s.shippingRepo.UpdateStatus(sh.ID, "cancelled", tenantID)
		}
	}

	// Publicar evento de cambio de estado
	if s.publisher != nil {
		ctx := context.Background()
		_ = s.publisher.PublishOrderStatusChanged(ctx, map[string]interface{}{
			"id":     orderID,
			"status": status,
		}, tenantID)

		// Si la orden fue pagada, publicar evento order.paid con los items para descontar stock
		// El descuento de stock se hará en product-micro al recibir este evento
		if status == "paid" {
			order, err := s.orderRepo.GetByID(orderID, tenantID)
			if err == nil && order != nil {
				// Preparar los items para el evento
				items := make([]map[string]interface{}, 0, len(order.Items))
				for _, item := range order.Items {
					items = append(items, map[string]interface{}{
						"product_id": item.ProductID,
						"quantity":   item.Quantity,
					})
				}

				_ = s.publisher.PublishOrderPaid(ctx, map[string]interface{}{
					"order_id": orderID,
					"items":    items,
				}, tenantID)
			}
		}
	}

	return nil
}

// generateTracking crea un tracking number simple
func generateTracking() string {
	return "TRK-" + uuid.New().String()[0:8]
}

func (s *Service) GetOrdersByProducer(producerID string, tenantID string) ([]dto.OrderDTO, error) {
	if producerID == "" || tenantID == "" {
		return nil, domain.BadRequestError{Message: "producerID y tenantID requeridos"}
	}

	producer, err := s.userRepo.GetByID(producerID, tenantID)
	if err != nil {
		return nil, err
	}
	fmt.Println(producer.Rol)
	if producer.Rol != "producer" && producer.Rol != "admin" {
		return nil, domain.UnauthorizedError{Message: "El usuario no es un productor válido"}
	}

	return s.orderRepo.GetOrdersByProducer(producerID, tenantID)
}
