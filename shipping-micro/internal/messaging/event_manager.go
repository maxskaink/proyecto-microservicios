package messaging

import (
	"context"
	"fmt"

	"github.com/maxskaink/proyecto-microservicios/shipping-micro/internal/db/repositories"
	"github.com/maxskaink/proyecto-microservicios/shipping-micro/internal/messaging/handlers"
	"github.com/maxskaink/proyecto-microservicios/shipping-micro/internal/messaging/rabbitmq"
	tenant "github.com/maxskaink/proyecto-microservicios/shipping-micro/internal/services/tenant"
	"github.com/maxskaink/proyecto-microservicios/shipping-micro/pkg/logger"
)

type EventManager struct{ consumer Consumer }

func NewEventManager(productRepo repositories.IProductRepository, userRepo repositories.IUserRepository, tenantService *tenant.TenantService) (*EventManager, error) {
	dispatcher := handlers.NewEventDispatcher(productRepo, userRepo, tenantService)
	cfg := rabbitmq.DefaultConfig()
	cm, err := rabbitmq.NewConnectionManager(cfg)
	if err != nil {
		return nil, fmt.Errorf("error al crear connection manager: %w", err)
	}
	consumer, err := rabbitmq.NewConsumer(cm, dispatcher)
	if err != nil {
		cm.Close()
		return nil, fmt.Errorf("error al crear consumidor: %w", err)
	}
	logger.Info("EventManager inicializado correctamente")
	return &EventManager{consumer: consumer}, nil
}

func (em *EventManager) Start(ctx context.Context) error {
	logger.Info("Iniciando consumo de eventos")
	return em.consumer.StartConsuming(ctx)
}
func (em *EventManager) Close() error {
	logger.Info("Cerrando gestor de eventos")
	return em.consumer.Close()
}
func (em *EventManager) GetConsumer() Consumer { return em.consumer }
