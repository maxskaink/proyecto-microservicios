package handlers

import (
	"fmt"

	"github.com/maxskaink/proyecto-microservicios/shipping-micro/pkg/logger"
)

type EventRegistry struct{ handlers map[string]EventHandler }

func NewEventRegistry() *EventRegistry {
	return &EventRegistry{handlers: make(map[string]EventHandler)}
}

func (er *EventRegistry) Register(handler EventHandler) error {
	et := handler.EventType()
	if _, ok := er.handlers[et]; ok {
		return fmt.Errorf("handler para evento '%s' ya registrado", et)
	}
	er.handlers[et] = handler
	logger.Info(fmt.Sprintf("Handler registrado para evento: %s", et))
	return nil
}

func (er *EventRegistry) Handle(eventType string, data []byte) error {
	fmt.Println("LLegó evento de tipo: " + eventType)
	h, ok := er.handlers[eventType]
	if !ok {
		logger.Info(fmt.Sprintf("No handler registrado para evento: %s", eventType))
		return nil
	}
	return h.Handle(data)
}

func (er *EventRegistry) GetHandler(eventType string) (EventHandler, error) {
	h, ok := er.handlers[eventType]
	if !ok {
		return nil, fmt.Errorf("handler no encontrado para evento: %s", eventType)
	}
	return h, nil
}

func (er *EventRegistry) ListHandlers() []string {
	out := make([]string, 0, len(er.handlers))
	for k := range er.handlers {
		out = append(out, k)
	}
	return out
}
