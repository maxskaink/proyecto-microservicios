package handlers

import (
	"fmt"

	"github.com/maxskaink/proyecto-microservicios/product-micro/internal/db/repositories"
	"github.com/maxskaink/proyecto-microservicios/product-micro/internal/messaging/events"
	"github.com/maxskaink/proyecto-microservicios/product-micro/pkg/logger"
)

// EventDispatcher mapea eventos a sus handlers
type EventDispatcher struct {
	handlers map[events.EventType]func([]byte) error
}

// NewEventDispatcher crea un nuevo dispatcher
func NewEventDispatcher(userRepository repositories.IUserRepository) *EventDispatcher {
	dispatcher := &EventDispatcher{
		handlers: make(map[events.EventType]func([]byte) error),
	}

	// Registrar handlers
	userHandler := &UserEventHandler{
		UserRepository: userRepository,
	}
	productHandler := &ProductEventHandler{}

	dispatcher.Register(events.UserCreated, userHandler.HandleUserCreated)
	dispatcher.Register(events.UserUpdated, userHandler.HandleUserUpdated)
	dispatcher.Register(events.UserDeleted, userHandler.HandleUserDeleted)
	dispatcher.Register(events.ProductCreated, productHandler.HandleProductCreated)
	dispatcher.Register(events.ProductUpdated, productHandler.HandleProductUpdated)

	return dispatcher
}

// Register registra un handler para un tipo de evento
func (ed *EventDispatcher) Register(eventType events.EventType, handler func([]byte) error) {
	ed.handlers[eventType] = handler
}

// Dispatch ejecuta el handler apropiado para un evento
func (ed *EventDispatcher) Dispatch(eventType events.EventType, data []byte) error {
	handler, exists := ed.handlers[eventType]
	if !exists {
		logger.Info(fmt.Sprintf("No handler registrado para el evento: %s", eventType))
		return nil // No es un error, simplemente ignoramos eventos desconocidos
	}

	return handler(data)
}
