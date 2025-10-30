package messaging

import (
	"github.com/maxskaink/proyecto-microservicios/product-micro/internal/messaging/rabbitmq"
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

// CreateConsumer crea un nuevo consumidor
func (f *Factory) CreateConsumer(queueName string, messageHandler func([]byte) error) (Consumer, error) {
	cm, err := rabbitmq.NewConnectionManager(f.config)
	if err != nil {
		return nil, err
	}

	consumer, err := rabbitmq.NewConsumer(cm, queueName, messageHandler)
	if err != nil {
		cm.Close()
		return nil, err
	}

	return consumer, nil
}
