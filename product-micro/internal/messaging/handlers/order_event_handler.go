package handlers

import (
	"encoding/json"
	"fmt"

	"github.com/maxskaink/proyecto-microservicios/product-micro/internal/db/repositories"
	"github.com/maxskaink/proyecto-microservicios/product-micro/internal/dto"
	"github.com/maxskaink/proyecto-microservicios/product-micro/pkg/logger"
)

// OrderPaidEvent representa la estructura del evento order.paid
type OrderPaidEvent struct {
	EventType string                 `json:"event_type"`
	Timestamp string                 `json:"timestamp"`
	TenantID  string                 `json:"tenant_id"`
	Data      map[string]interface{} `json:"data"`
}

// OrderPaidPayload representa el payload del evento order.paid
type OrderPaidPayload struct {
	OrderID string      `json:"order_id"`
	Items   []OrderItem `json:"items"`
}

// OrderItem representa un item de la orden
type OrderItem struct {
	ProductID string `json:"product_id"`
	Quantity  int    `json:"quantity"`
}

// OrderPaidHandler maneja eventos de orden pagada para descontar stock
type OrderPaidHandler struct {
	BaseHandler
	productRepo repositories.IProductRepository
}

// NewOrderPaidHandler crea un nuevo handler de orden pagada
func NewOrderPaidHandler(productRepo repositories.IProductRepository) *OrderPaidHandler {
	return &OrderPaidHandler{
		BaseHandler: NewBaseHandler("order.paid"),
		productRepo: productRepo,
	}
}

func (h *OrderPaidHandler) Handle(data []byte) error {
	var event OrderPaidEvent
	if err := json.Unmarshal(data, &event); err != nil {
		logger.Error(fmt.Sprintf("Error al parsear evento order.paid: %v", err))
		return err
	}

	logger.Info(fmt.Sprintf("Evento order.paid recibido - EventType: %s, TenantID: %s",
		event.EventType, event.TenantID))

	tenantID := event.TenantID
	if tenantID == "" {
		logger.Error("TenantID vacío en evento order.paid")
		return fmt.Errorf("tenant_id es requerido")
	}

	// Extraer el payload
	payload, ok := event.Data["payload"]
	if !ok {
		logger.Error("payload no encontrado en evento order.paid")
		return fmt.Errorf("payload no encontrado")
	}

	// Convertir payload a estructura
	payloadBytes, err := json.Marshal(payload)
	if err != nil {
		return err
	}

	var orderPaid OrderPaidPayload
	if err := json.Unmarshal(payloadBytes, &orderPaid); err != nil {
		logger.Error(fmt.Sprintf("Error al parsear payload de order.paid: %v. Payload raw: %s",
			err, string(payloadBytes)))
		return err
	}

	// Validar que tengamos items
	if len(orderPaid.Items) == 0 {
		logger.Error(fmt.Sprintf("Orden %s no tiene items para procesar", orderPaid.OrderID))
		return fmt.Errorf("orden sin items")
	}

	logger.Info(fmt.Sprintf("Procesando orden pagada %s para tenant %s con %d items",
		orderPaid.OrderID, tenantID, len(orderPaid.Items)))

	// Descontar stock para cada item de la orden
	successCount := 0
	for i, item := range orderPaid.Items {
		logger.Info(fmt.Sprintf("Procesando item %d/%d - ProductID: %s, Quantity: %d",
			i+1, len(orderPaid.Items), item.ProductID, item.Quantity))

		// Validar datos del item
		if item.ProductID == "" {
			logger.Error(fmt.Sprintf("Item %d tiene ProductID vacío, saltando...", i+1))
			continue
		}
		if item.Quantity <= 0 {
			logger.Error(fmt.Sprintf("Item %d tiene cantidad inválida (%d), saltando...", i+1, item.Quantity))
			continue
		}

		// Obtener el producto actual
		product, err := h.productRepo.GetByIdProduct(item.ProductID, tenantID)
		if err != nil {
			logger.Error(fmt.Sprintf("Error al obtener producto %s: %v", item.ProductID, err))
			continue // Continuar con el siguiente item aunque uno falle
		}

		if product == nil {
			logger.Error(fmt.Sprintf("Producto %s no encontrado en tenant %s", item.ProductID, tenantID))
			continue
		}

		// Verificar si hay suficiente stock
		if product.Stock < item.Quantity {
			logger.Error(fmt.Sprintf("Stock insuficiente para producto %s. Stock actual: %d, requerido: %d",
				item.ProductID, product.Stock, item.Quantity))
			// Continuar con el siguiente item - en producción podrías querer manejar esto de otra forma
			continue
		}

		// Descontar el stock
		newStock := product.Stock - item.Quantity
		logger.Info(fmt.Sprintf("Descontando stock del producto %s: %d -> %d",
			item.ProductID, product.Stock, newStock))

		// Actualizar el producto
		product.Stock = newStock
		productRequest := dto.ProductDTOResponseTORequest(*product)

		_, err = h.productRepo.UpdateProduct(item.ProductID, &productRequest, tenantID)
		if err != nil {
			logger.Error(fmt.Sprintf("Error al actualizar stock del producto %s: %v", item.ProductID, err))
			continue
		}

		successCount++
		logger.Info(fmt.Sprintf("Stock actualizado exitosamente para producto %s (nuevo stock: %d)",
			item.ProductID, newStock))
	}

	logger.Info(fmt.Sprintf("Orden pagada %s procesada: %d/%d productos actualizados exitosamente",
		orderPaid.OrderID, successCount, len(orderPaid.Items)))
	return nil
}
