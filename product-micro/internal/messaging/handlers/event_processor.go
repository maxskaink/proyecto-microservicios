package handlers

import (
	"encoding/json"
	"fmt"

	"github.com/maxskaink/proyecto-microservicios/product-micro/internal/db/repositories"
	"github.com/maxskaink/proyecto-microservicios/product-micro/internal/messaging/events"
	"github.com/maxskaink/proyecto-microservicios/product-micro/pkg/logger"
)

// EventProcessor procesa mensajes y los envía al dispatcher correcto
type EventProcessor struct {
	dispatcher *EventDispatcher
}

// NewEventProcessor crea un nuevo procesador de eventos
func NewEventProcessor(userRepository repositories.IUserRepository) *EventProcessor {
	return &EventProcessor{
		dispatcher: NewEventDispatcher(userRepository),
	}
}

// ProcessMessage procesa un mensaje y lo envía al handler correspondiente
func (ep *EventProcessor) ProcessMessage(data []byte, routingKey string) error {
	// Validar que el mensaje no esté vacío
	if len(data) == 0 {
		logger.Error("Mensaje vacío recibido")
		return fmt.Errorf("mensaje vacío")
	}

	logger.Info(fmt.Sprintf("EventProcessor: procesando mensaje con routing key: %s", routingKey))

	// Parsear el evento para obtener su tipo
	var baseEvent events.Event
	if err := json.Unmarshal(data, &baseEvent); err != nil {
		logger.Error(fmt.Sprintf("EventProcessor: error al parsear evento base: %v", err))
		// Retornar error para que sea descartado (no reintentar)
		return fmt.Errorf("mensaje JSON inválido: %w", err)
	}

	// Determinar el tipo de evento
	eventType := baseEvent.Type
	if eventType == "" {
		eventType = events.EventType(routingKey)
	}

	// Validar que sea un tipo conocido
	if !isValidEventType(eventType) {
		logger.Error(fmt.Sprintf("EventProcessor: tipo de evento desconocido: %s (routing key: %s)",
			eventType, routingKey))
		// Descartar mensajes con tipos desconocidos
		return fmt.Errorf("tipo de evento desconocido: %s", eventType)
	}

	logger.Info(fmt.Sprintf("EventProcessor: enviando evento tipo %s al dispatcher", eventType))

	// Dispatcher maneja el evento basado en su tipo
	if err := ep.dispatcher.Dispatch(eventType, data); err != nil {
		logger.Error(fmt.Sprintf("EventProcessor: error en dispatcher para evento %s: %v",
			eventType, err))
		return err
	}

	logger.Info(fmt.Sprintf("EventProcessor: evento %s procesado exitosamente", eventType))
	return nil
}

// isValidEventType valida que el tipo de evento sea conocido
func isValidEventType(eventType events.EventType) bool {
	validTypes := map[events.EventType]bool{
		events.UserCreated:    true,
		events.UserUpdated:    true,
		events.UserDeleted:    true,
		events.ProductCreated: true,
		events.ProductUpdated: true,
	}
	return validTypes[eventType]
}
