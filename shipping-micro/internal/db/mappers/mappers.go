package mappers

import (
	"github.com/maxskaink/proyecto-microservicios/shipping-micro/internal/db/models"
	"github.com/maxskaink/proyecto-microservicios/shipping-micro/internal/dto"
)

// ProductDBToDTO convierte ProductDB a ProductDTO
func ProductDBToDTO(p *models.ProductDB) *dto.ProductDTO {
	if p == nil {
		return nil
	}
	return &dto.ProductDTO{
		ID:          p.ID,
		Name:        p.Name,
		Description: p.Description,
		Price:       p.Price,
		Stock:       p.Stock,
		CreatedAt:   p.CreatedAt,
		UpdatedAt:   p.UpdatedAt,
	}
}

// UserDBToDTO convierte UserDB a UserDTO
func UserDBToDTO(u *models.UserDB) *dto.UserDTO {
	if u == nil {
		return nil
	}
	return &dto.UserDTO{
		ID:        u.ID,
		Email:     u.Email,
		Name:      u.Name,
		CreatedAt: u.CreatedAt,
		UpdatedAt: u.UpdatedAt,
	}
}

// CartItemDBToDTO convierte CartItemDB a CartItemDTO
func CartItemDBToDTO(c *models.CartItemDB) *dto.CartItemDTO {
	if c == nil {
		return nil
	}
	return &dto.CartItemDTO{
		ID:        c.ID,
		UserID:    c.UserID,
		ProductID: c.ProductID,
		Quantity:  c.Quantity,
		CreatedAt: c.CreatedAt,
		UpdatedAt: c.UpdatedAt,
		Product:   ProductDBToDTO(&c.Product),
	}
}

// OrderDBToDTO convierte OrderDB a OrderDTO
func OrderDBToDTO(o *models.OrderDB) *dto.OrderDTO {
	if o == nil {
		return nil
	}

	var items []dto.OrderItemDTO
	for _, item := range o.Items {
		items = append(items, *OrderItemDBToDTO(&item))
	}

	return &dto.OrderDTO{
		ID:         o.ID,
		UserID:     o.UserID,
		TotalPrice: o.TotalPrice,
		Status:     o.Status,
		Items:      items,
		CreatedAt:  o.CreatedAt,
		UpdatedAt:  o.UpdatedAt,
	}
}

// OrderItemDBToDTO convierte OrderItemDB a OrderItemDTO
func OrderItemDBToDTO(oi *models.OrderItemDB) *dto.OrderItemDTO {
	if oi == nil {
		return nil
	}
	return &dto.OrderItemDTO{
		ID:        oi.ID,
		OrderID:   oi.OrderID,
		ProductID: oi.ProductID,
		Quantity:  oi.Quantity,
		Price:     oi.Price,
		CreatedAt: oi.CreatedAt,
	}
}

// ShippingDBToDTO convierte ShippingDB a ShippingDTO
func ShippingDBToDTO(s *models.ShippingDB) *dto.ShippingDTO {
	if s == nil {
		return nil
	}
	return &dto.ShippingDTO{
		ID:              s.ID,
		OrderID:         s.OrderID,
		TrackingNumber:  s.TrackingNumber,
		ShippingAddress: s.ShippingAddress,
		Status:          s.Status,
		CreatedAt:       s.CreatedAt,
		UpdatedAt:       s.UpdatedAt,
	}
}
