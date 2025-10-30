package rabbitmq

import (
	"context"
	"fmt"

	"github.com/maxskaink/proyecto-microservicios/product-micro/pkg/logger"
	amqp "github.com/rabbitmq/amqp091-go"
)

const (
	// Constantes para la configuración de RabbitMQ
	ExchangeName   = "user_events"
	ExchangeType   = "topic"
	UserCreatedKey = "user.created"
	UserUpdatedKey = "user.updated"
	UserDeletedKey = "user.deleted"
	QueueName      = "product_user_events"
)

// Consumer maneja el consumo de mensajes
type Consumer struct {
	cm             *ConnectionManager
	messageHandler func([]byte) error
	queueName      string
}

// NewConsumer crea un nuevo consumidor
func NewConsumer(cm *ConnectionManager, queueName string, messageHandler func([]byte) error) (*Consumer, error) {
	c := &Consumer{
		cm:             cm,
		messageHandler: messageHandler,
		queueName:      queueName,
	}

	// Inicializar el exchange y la cola
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

	// Declarar el exchange
	err = ch.ExchangeDeclare(
		ExchangeName, // nombre
		ExchangeType, // tipo
		true,         // durable
		false,        // auto-eliminado
		false,        // interno
		false,        // no-wait
		nil,          // argumentos
	)
	if err != nil {
		return err
	}

	// Declarar la cola
	q, err := ch.QueueDeclare(
		c.queueName, // nombre
		true,        // durable
		false,       // auto-eliminado
		false,       // exclusiva
		false,       // no-wait
		nil,         // argumentos
	)
	if err != nil {
		return err
	}

	// Vincular la cola al exchange
	err = ch.QueueBind(
		q.Name,       // nombre de la cola
		"user.#",     // routing key (todos los eventos de usuarios)
		ExchangeName, // nombre del exchange
		false,        // no-wait
		nil,          // argumentos
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

	// Configurar QoS
	err = ch.Qos(
		1,     // prefetch count (procesar un mensaje a la vez)
		0,     // prefetch size
		false, // global
	)
	if err != nil {
		return err
	}

	// Consumir mensajes
	msgs, err := ch.Consume(
		c.queueName, // cola
		"",          // consumer
		false,       // auto-ack
		false,       // exclusive
		false,       // no-local
		false,       // no-wait
		nil,         // args
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

				if err := c.processMessage(msg); err != nil {
					logger.Error(fmt.Sprintf("Error al procesar mensaje: %v", err))
					msg.Nack(false, true)
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
	fmt.Print("Mensaje recibido")
	logger.Info(fmt.Sprintf("Mensaje recibido: routing key = %s, body = %s", msg.RoutingKey, string(msg.Body)))

	// Llamar al manejador de mensajes
	if c.messageHandler != nil {
		return c.messageHandler(msg.Body)
	}

	return nil
}

// Close cierra el consumidor
func (c *Consumer) Close() error {
	return c.cm.Close()
}
