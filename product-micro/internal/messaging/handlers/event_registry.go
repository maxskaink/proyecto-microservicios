package handlers

import (
	"fmt"

	"github.com/maxskaink/proyecto-microservicios/product-micro/pkg/logger"
)

// EventRegistry mantiene un registro de todos los handlers por tipo de evento
type EventRegistry struct {
	handlers map[string]EventHandler
}

// NewEventRegistry crea un nuevo registro de eventos
func NewEventRegistry() *EventRegistry {
	return &EventRegistry{
		handlers: make(map[string]EventHandler),
	}
}

// Register registra un handler para un tipo de evento
func (er *EventRegistry) Register(handler EventHandler) error {
	eventType := handler.EventType()

	if _, exists := er.handlers[eventType]; exists {
		return fmt.Errorf("handler para evento '%s' ya está registrado", eventType)
	}

	er.handlers[eventType] = handler
	logger.Info(fmt.Sprintf("Handler registrado para evento: %s", eventType))

	return nil
}

// Handle procesa un evento basado en su tipo
func (er *EventRegistry) Handle(eventType string, data []byte) error {
	handler, exists := er.handlers[eventType]
	if !exists {
		logger.Info(fmt.Sprintf("No handler registrado para evento: %s", eventType))
		return nil // No es un error, simplemente ignoramos eventos desconocidos
	}

	return handler.Handle(data)
}

// GetHandler obtiene un handler específico
func (er *EventRegistry) GetHandler(eventType string) (EventHandler, error) {
	handler, exists := er.handlers[eventType]
	if !exists {
		return nil, fmt.Errorf("handler no encontrado para evento: %s", eventType)
	}
	return handler, nil
}

// ListHandlers retorna todos los tipos de eventos registrados
func (er *EventRegistry) ListHandlers() []string {
	var eventTypes []string
	for eventType := range er.handlers {
		eventTypes = append(eventTypes, eventType)
	}
	return eventTypes
}
