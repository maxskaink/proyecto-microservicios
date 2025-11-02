package handlers

import (
	"encoding/json"
	"fmt"

	"github.com/maxskaink/proyecto-microservicios/product-micro/pkg/logger"
)

// EventHandler define la interfaz base para todos los handlers
type EventHandler interface {
	// Handle procesa el evento
	Handle(data []byte) error
	// EventType retorna el tipo de evento que maneja
	EventType() string
}

// BaseHandler proporciona funcionalidad común para todos los handlers
type BaseHandler struct {
	eventType string
}

// NewBaseHandler crea un nuevo BaseHandler
func NewBaseHandler(eventType string) BaseHandler {
	return BaseHandler{eventType: eventType}
}

// EventType retorna el tipo de evento
func (b *BaseHandler) EventType() string {
	return b.eventType
}

// UnmarshalEvent desmarshala los datos del evento de forma segura
func (b *BaseHandler) UnmarshalEvent(data []byte, v interface{}) error {
	if err := json.Unmarshal(data, v); err != nil {
		logger.Error(fmt.Sprintf("Error al deserializar evento %s: %v", b.eventType, err))
		return err
	}
	return nil
}

// LogEvent registra el procesamiento de un evento
func (b *BaseHandler) LogEvent(action string) {
	logger.Info(fmt.Sprintf("[%s] %s", b.eventType, action))
}
