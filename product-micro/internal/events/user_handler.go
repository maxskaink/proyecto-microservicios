package events

import (
	"context"
	"fmt"

	"github.com/maxskaink/proyecto-microservicios/product-micro/internal/db/repositories"
	"github.com/maxskaink/proyecto-microservicios/product-micro/internal/messaging"
	"github.com/maxskaink/proyecto-microservicios/product-micro/internal/messaging/handlers"
	"github.com/maxskaink/proyecto-microservicios/product-micro/internal/messaging/rabbitmq"
	"github.com/maxskaink/proyecto-microservicios/product-micro/internal/services/tenant"
	"github.com/maxskaink/proyecto-microservicios/product-micro/pkg/logger"
)

// UserHandler maneja los eventos de usuarios
type UserHandler struct {
	consumer messaging.Consumer
}

// NewUserHandler crea un nuevo manejador de eventos de usuarios
func NewUserHandler(userRepository repositories.IUserRepository, tenantService *tenant.TenantService) (*UserHandler, error) {
	// Crear el dispatcher con todos los handlers registrados
	dispatcher := handlers.NewEventDispatcher(userRepository, tenantService)

	// Crear el connection manager
	config := rabbitmq.DefaultConfig()
	cm, err := rabbitmq.NewConnectionManager(config)
	if err != nil {
		return nil, fmt.Errorf("error al crear connection manager: %w", err)
	}

	// Crear consumidor pasando el dispatcher como procesador de eventos
	consumer, err := rabbitmq.NewConsumer(cm, dispatcher)
	if err != nil {
		cm.Close()
		return nil, fmt.Errorf("error al crear consumidor de eventos: %w", err)
	}

	logger.Info("UserHandler inicializado correctamente con EventDispatcher")
	return &UserHandler{
		consumer: consumer,
	}, nil
}

// Start inicia el consumo de mensajes
func (h *UserHandler) Start(ctx context.Context) error {
	logger.Info("Iniciando consumo de eventos de usuarios")
	return h.consumer.StartConsuming(ctx)
}

// Close cierra el consumidor
func (h *UserHandler) Close() error {
	logger.Info("Cerrando consumidor de eventos")
	return h.consumer.Close()
}
