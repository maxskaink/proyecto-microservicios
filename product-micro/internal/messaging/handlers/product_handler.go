package handlers

import (
	"encoding/json"
	"fmt"

	"github.com/maxskaink/proyecto-microservicios/product-micro/internal/messaging/events"
	"github.com/maxskaink/proyecto-microservicios/product-micro/pkg/logger"
)

// ProductEventHandler maneja eventos de producto
type ProductEventHandler struct{}

// HandleProductCreated procesa el evento de producto creado
func (h *ProductEventHandler) HandleProductCreated(data []byte) error {
	var event events.ProductCreatedEvent

	if err := json.Unmarshal(data, &event); err != nil {
		logger.Error(fmt.Sprintf("Error al deserializar evento de producto creado: %v", err))
		return err
	}

	logger.Info(fmt.Sprintf("Producto creado: ID=%s, Productor=%s",
		event.ProductID, event.ProducerID))

	// Lógica adicional si es necesario
	// Por ejemplo: indexar en búsqueda, notificar a otros servicios, etc.

	return nil
}

// HandleProductUpdated procesa el evento de producto actualizado
func (h *ProductEventHandler) HandleProductUpdated(data []byte) error {
	var event events.ProductUpdatedEvent

	if err := json.Unmarshal(data, &event); err != nil {
		logger.Error(fmt.Sprintf("Error al deserializar evento de producto actualizado: %v", err))
		return err
	}

	logger.Info(fmt.Sprintf("Producto actualizado: ID=%s, Productor=%s",
		event.ProductID, event.ProducerID))

	// Lógica adicional si es necesario

	return nil
}
