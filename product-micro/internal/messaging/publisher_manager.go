package messaging

import "github.com/maxskaink/proyecto-microservicios/product-micro/internal/messaging/rabbitmq"

// NewProductEventPublisher crea un Publisher para eventos de productos usando RabbitMQ
func NewProductEventPublisher() (Publisher, error) {
	config := rabbitmq.DefaultConfig()
	return rabbitmq.NewProductPublisher(config)
}
