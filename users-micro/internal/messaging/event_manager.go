package messaging

import (
	"context"
	"fmt"
	"os"

	"github.com/maxskaink/proyecto-microservicios/users-micro/internal/messaging/rabbitmq"
	"github.com/maxskaink/proyecto-microservicios/users-micro/pkg/logger"
)

// Consumer define la interfaz para el consumo de mensajes
type Consumer interface {
	StartConsuming(ctx context.Context) error
	Close() error
}

// Dispatcher define la interfaz para procesar mensajes
type Dispatcher interface {
	ProcessMessage(data []byte, routingKey string) error
	Close() error
}

// EventManager gestiona el consumo de eventos de RabbitMQ
// Inicializa el dispatcher y el consumer para procesar eventos de tenant
type EventManager struct {
	consumer Consumer
}

// NewEventManager crea un nuevo gestor de eventos
// Recibe el dispatcher que contiene todos los handlers y crea el consumer de RabbitMQ
func NewEventManager(dispatcher Dispatcher) (*EventManager, error) {
	// Obtener URL de RabbitMQ
	rabbitmqURL := os.Getenv("RABBITMQ_URL")
	if rabbitmqURL == "" {
		rabbitmqURL = "amqp://guest:guest@localhost:5672/"
	}

	// Crear consumer de RabbitMQ que escucha eventos de tenant
	consumer, err := rabbitmq.NewConsumer(rabbitmqURL, dispatcher)
	if err != nil {
		return nil, fmt.Errorf("error al crear consumidor de eventos: %w", err)
	}

	logger.Info("EventManager inicializado correctamente con consumer y dispatcher de eventos")
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
	if em.consumer != nil {
		return em.consumer.Close()
	}
	return nil
}
