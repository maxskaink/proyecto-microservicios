package messaging

import (
	"github.com/maxskaink/proyecto-microservicios/users-micro/internal/messaging/rabbitmq"
)

// Factory crea instancias de los componentes de mensajería
type Factory struct {
	config *rabbitmq.Config
}

// NewFactory crea una nueva fábrica de componentes de mensajería
func NewFactory(config *rabbitmq.Config) *Factory {
	if config == nil {
		config = rabbitmq.DefaultConfig()
	}
	return &Factory{
		config: config,
	}
}

// CreatePublisher crea un nuevo publicador
func (f *Factory) CreatePublisher() (Publisher, error) {
	cm, err := rabbitmq.NewConnectionManager(f.config)
	if err != nil {
		return nil, err
	}

	publisher, err := rabbitmq.NewPublisher(cm)
	if err != nil {
		cm.Close()
		return nil, err
	}

	return publisher, nil
}
