package rabbitmq

import (
	"context"
	"fmt"

	"github.com/maxskaink/proyecto-microservicios/product-micro/pkg/logger"
	amqp "github.com/rabbitmq/amqp091-go"
)

const (
	ExchangeName   = "user_events"
	ExchangeType   = "topic"
	UserCreatedKey = "user.created"
	UserUpdatedKey = "user.updated"
	UserDeletedKey = "user.deleted"
	QueueName      = "product_user_events"
)

type Consumer struct {
	cm         *ConnectionManager
	queueName  string
	dispatcher Dispatcher
}

// Dispatcher define la interfaz para procesar mensajes
type Dispatcher interface {
	ProcessMessage(data []byte, routingKey string) error
}

// NewConsumer crea un nuevo consumidor
func NewConsumer(cm *ConnectionManager, queueName string, dispatcher Dispatcher) (*Consumer, error) {
	c := &Consumer{
		cm:         cm,
		queueName:  queueName,
		dispatcher: dispatcher,
	}

	if err := c.setupQueue(); err != nil {
		return nil, err
	}

	return c, nil
}

// setupQueue inicializa el exchange y la cola para consumir mensajes
func (c *Consumer) setupQueue() error {
	ch, err := c.cm.Channel()
	if err != nil {
		return err
	}
	defer ch.Close()

	err = ch.ExchangeDeclare(
		ExchangeName,
		ExchangeType,
		true,  // durable
		false, // auto-eliminated
		false, // internal
		false, // no-wait
		nil,   // args
	)
	if err != nil {
		return err
	}

	q, err := ch.QueueDeclare(
		c.queueName,
		true,  // durable
		false, // auto-eliminated
		false, // exclusive
		false, // no-wait
		nil,   // args
	)
	if err != nil {
		return err
	}

	err = ch.QueueBind(
		q.Name,
		"user.#",
		ExchangeName,
		false,
		nil,
	)

	return err
}

// StartConsuming comienza a consumir mensajes
func (c *Consumer) StartConsuming(ctx context.Context) error {
	ch, err := c.cm.Channel()
	if err != nil {
		return err
	}
	defer ch.Close()

	err = ch.Qos(1, 0, false)
	if err != nil {
		return err
	}

	msgs, err := ch.Consume(
		c.queueName,
		"",    // consumer
		false, // auto-ack
		false, // exclusive
		false, // no-local
		false, // no-wait
		nil,   // args
	)
	if err != nil {
		return err
	}

	logger.Info(fmt.Sprintf("Comenzando a consumir mensajes de la cola %s", c.queueName))
	done := make(chan struct{})

	go func() {
		for {
			select {
			case <-ctx.Done():
				logger.Info("Deteniendo consumo de mensajes")
				close(done)
				return

			case msg, ok := <-msgs:
				if !ok {
					logger.Error("Canal de mensajes cerrado")
					close(done)
					return
				}

				// Procesar mensaje con manejo de errores mejorado
				if err := c.processMessage(msg); err != nil {
					logger.Error(fmt.Sprintf("Error al procesar mensaje (routing key: %s): %v", msg.RoutingKey, err))
					// NO reintentar: descartar el mensaje (dead letter)
					msg.Ack(false)
				} else {
					msg.Ack(false)
				}
			}
		}
	}()

	<-done
	return nil
}

// processMessage procesa un mensaje recibido
func (c *Consumer) processMessage(msg amqp.Delivery) error {
	logger.Info(fmt.Sprintf("Procesando mensaje - Routing key: %s, Body: %s",
		msg.RoutingKey, string(msg.Body)))

	if c.dispatcher == nil {
		return fmt.Errorf("dispatcher no configurado")
	}

	return c.dispatcher.ProcessMessage(msg.Body, msg.RoutingKey)
}

// Close cierra el consumidor
func (c *Consumer) Close() error {
	return c.cm.Close()
}
