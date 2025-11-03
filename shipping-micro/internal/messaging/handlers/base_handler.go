package handlers

import (
	"encoding/json"
	"fmt"

	"github.com/maxskaink/proyecto-microservicios/shipping-micro/pkg/logger"
)

// EventHandler interfaz base de handlers
type EventHandler interface {
	Handle(data []byte) error
	EventType() string
}

type BaseHandler struct{ eventType string }

func NewBaseHandler(eventType string) BaseHandler { return BaseHandler{eventType: eventType} }
func (b *BaseHandler) EventType() string          { return b.eventType }
func (b *BaseHandler) UnmarshalEvent(data []byte, v interface{}) error {
	if err := json.Unmarshal(data, v); err != nil {
		logger.Error(fmt.Sprintf("Error al deserializar evento %s: %v", b.eventType, err))
		return err
	}
	return nil
}
func (b *BaseHandler) LogEvent(action string) {
	logger.Info(fmt.Sprintf("[%s] %s", b.eventType, action))
}
