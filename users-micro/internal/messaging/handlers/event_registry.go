package handlers

import (
	"fmt"

	"github.com/maxskaink/proyecto-microservicios/users-micro/pkg/logger"
)

// EventRegistry mantiene un registro centralizado de handlers de eventos
type EventRegistry struct {
	handlers map[string]EventHandler
}

// NewEventRegistry crea un nuevo registro de handlers
func NewEventRegistry() *EventRegistry {
	return &EventRegistry{
		handlers: make(map[string]EventHandler),
	}
}

// Register registra un nuevo handler para un tipo de evento
func (r *EventRegistry) Register(handler EventHandler) error {
	eventType := handler.EventType()
	if _, exists := r.handlers[eventType]; exists {
		logger.Error(fmt.Sprintf("Handler para evento %s ya está registrado", eventType))
		return fmt.Errorf("handler for event type %s already registered", eventType)
	}
	r.handlers[eventType] = handler
	logger.Info(fmt.Sprintf("Handler registrado para evento: %s", eventType))
	return nil
}

// Handle procesa un evento buscando el handler apropiado
func (r *EventRegistry) Handle(eventType string, data []byte) error {
	handler, exists := r.handlers[eventType]
	if !exists {
		logger.Error(fmt.Sprintf("No hay handler registrado para evento: %s", eventType))
		return fmt.Errorf("no handler for event type: %s", eventType)
	}

	return handler.Handle(data)
}

// GetHandler retorna el handler para un tipo de evento específico
func (r *EventRegistry) GetHandler(eventType string) EventHandler {
	return r.handlers[eventType]
}

// ListHandlers retorna una lista de todos los tipos de eventos registrados
func (r *EventRegistry) ListHandlers() []string {
	var handlers []string
	for eventType := range r.handlers {
		handlers = append(handlers, eventType)
	}
	return handlers
}
