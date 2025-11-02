package messaging

import (
	"github.com/maxskaink/proyecto-microservicios/product-micro/internal/messaging/rabbitmq"
)

// Factory crea instancias de los componentes de mensajería
type Factory struct {
	config *rabbitmq.Config
}

// messageHandlerAdapter convierte una función en un Dispatcher
type messageHandlerAdapter struct {
	handler func([]byte) error
}

// ProcessMessage implementa la interfaz Dispatcher
func (a *messageHandlerAdapter) ProcessMessage(data []byte, routingKey string) error {
	return a.handler(data)
}

func (f *Factory) CreateConsumer(queueName string, messageHandler func([]byte) error) (Consumer, error) {
	cm, err := rabbitmq.NewConnectionManager(f.config)
	if err != nil {
		return nil, err
	}

	// Adaptar la función a Dispatcher
	dispatcher := &messageHandlerAdapter{handler: messageHandler}

	consumer, err := rabbitmq.NewConsumer(cm, dispatcher)
	if err != nil {
		cm.Close()
		return nil, err
	}

	return consumer, nil
}
