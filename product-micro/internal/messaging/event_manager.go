package messaging

import (
	"context"
	"fmt"

	"github.com/maxskaink/proyecto-microservicios/product-micro/internal/db/repositories"
	"github.com/maxskaink/proyecto-microservicios/product-micro/internal/messaging/handlers"
	"github.com/maxskaink/proyecto-microservicios/product-micro/internal/messaging/rabbitmq"
	"github.com/maxskaink/proyecto-microservicios/product-micro/internal/services/tenant"
	"github.com/maxskaink/proyecto-microservicios/product-micro/pkg/logger"
)

// EventManager gestiona el consumo de eventos de RabbitMQ
// Inicializa el dispatcher y el consumer para procesar todos los eventos
type EventManager struct {
	consumer Consumer
}

// NewEventManager crea un nuevo gestor de eventos
// Configura el dispatcher con todos los handlers y el consumer de RabbitMQ
func NewEventManager(userRepository repositories.IUserRepository, productRepository repositories.IProductRepository, tenantService *tenant.TenantService) (*EventManager, error) {
	// Crear el dispatcher con todos los handlers registrados
	// (user.created, user.updated, user.deleted, tenant.created, tenant.deleted, order.paid)
	dispatcher := handlers.NewEventDispatcher(userRepository, productRepository, tenantService)

	// Crear el connection manager para conectarse a RabbitMQ
	config := rabbitmq.DefaultConfig()
	cm, err := rabbitmq.NewConnectionManager(config)
	if err != nil {
		return nil, fmt.Errorf("error al crear connection manager: %w", err)
	}

	// Crear consumidor que escucha eventos de los exchanges configurados
	consumer, err := rabbitmq.NewConsumer(cm, dispatcher)
	if err != nil {
		cm.Close()
		return nil, fmt.Errorf("error al crear consumidor de eventos: %w", err)
	}

	logger.Info("EventManager inicializado correctamente con todos los handlers")
	return &EventManager{
		consumer: consumer,
	}, nil
}

// Start inicia el consumo de mensajes desde RabbitMQ
// Esta es una operación bloqueante, típicamente se ejecuta en una goroutine
func (em *EventManager) Start(ctx context.Context) error {
	logger.Info("Iniciando consumo de eventos desde RabbitMQ")
	return em.consumer.StartConsuming(ctx)
}

// Close cierra la conexión del consumidor a RabbitMQ
func (em *EventManager) Close() error {
	logger.Info("Cerrando gestor de eventos")
	return em.consumer.Close()
}

// GetDispatcher retorna el dispatcher para acceso directo si es necesario
// (utilidad para testing o casos especiales)
func (em *EventManager) GetConsumer() Consumer {
	return em.consumer
}
