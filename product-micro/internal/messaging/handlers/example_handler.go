package handlers

import (
	"encoding/json"
	"fmt"

	"github.com/maxskaink/proyecto-microservicios/product-micro/pkg/logger"
)

// EJEMPLO: Esto muestra cómo agregar un nuevo tipo de evento en el futuro
// sin cambiar el consumer ni el dispatcher
//
// Pasos para agregar un nuevo evento:
// 1. Definir el tipo en events/event.go
// 2. Crear un handler como este
// 3. Registrarlo en NewEventDispatcher() en dispatcher.go

// OrderEventHandler sería un handler de ejemplo para eventos de órdenes
type OrderEventHandler struct{}

// HandleOrderCreated procesaría eventos de órdenes creadas
func (h *OrderEventHandler) HandleOrderCreated(data []byte) error {
	type OrderEvent struct {
		OrderID string  `json:"order_id"`
		UserID  string  `json:"user_id"`
		Total   float64 `json:"total"`
	}

	var order OrderEvent
	if err := json.Unmarshal(data, &order); err != nil {
		logger.Error(fmt.Sprintf("Error al deserializar evento de orden: %v", err))
		return err
	}

	logger.Info(fmt.Sprintf("Orden creada: ID=%s, Usuario=%s, Total=%f",
		order.OrderID, order.UserID, order.Total))

	// Tu lógica aquí...
	return nil
}

// NOTA: Descomentar estas líneas en dispatcher.go cuando necesites agregar el evento:
//
// orderHandler := &OrderEventHandler{}
// dispatcher.Register(events.OrderCreated, orderHandler.HandleOrderCreated)
