package rabbitmq

import (
	"context"
	"fmt"

	"github.com/maxskaink/proyecto-microservicios/users-micro/pkg/logger"
	amqp "github.com/rabbitmq/amqp091-go"
)

// Tenant Events
const (
	TenantExchangeName = "tenants_events"
	TenantExchangeType = "topic"
	TenantCreatedKey   = "tenant.created"
	TenantDeletedKey   = "tenant.deleted"
	TenantQueueName    = "users_tenant_events"
)

type Consumer struct {
	conn       *amqp.Connection
	url        string
	dispatcher Dispatcher
}

// Dispatcher define la interfaz para procesar mensajes
type Dispatcher interface {
	ProcessMessage(data []byte, routingKey string) error
}

// NewConsumer crea un nuevo consumidor de eventos
func NewConsumer(url string, dispatcher Dispatcher) (*Consumer, error) {
	consumer := &Consumer{
		url:        url,
		dispatcher: dispatcher,
	}

	// Conectar a RabbitMQ
	if err := consumer.connect(); err != nil {
		return nil, fmt.Errorf("error al conectar a RabbitMQ: %w", err)
	}

	// Configurar exchanges y colas
	if err := consumer.setupTenantQueue(); err != nil {
		consumer.Close()
		return nil, fmt.Errorf("error al configurar cola de tenants: %w", err)
	}

	return consumer, nil
}

// connect establece la conexión a RabbitMQ
func (c *Consumer) connect() error {
	var err error
	// Intentar conectar con reintentos
	for i := 0; i < 5; i++ {
		c.conn, err = amqp.Dial(c.url)
		if err == nil {
			logger.Info("Conectado a RabbitMQ para consumidor de eventos")
			return nil
		}
		logger.Error(fmt.Sprintf("Intento %d de conexión a RabbitMQ fallido: %v", i+1, err))
	}
	return err
}

// setupTenantQueue configura el exchange y la cola para eventos de tenants
func (c *Consumer) setupTenantQueue() error {
	ch, err := c.conn.Channel()
	if err != nil {
		return err
	}
	defer ch.Close()

	// Declarar exchange
	err = ch.ExchangeDeclare(
		TenantExchangeName,
		TenantExchangeType,
		true,  // durable
		false, // auto-eliminated
		false, // internal
		false, // no-wait
		nil,   // args
	)
	if err != nil {
		return err
	}

	// Declarar cola
	q, err := ch.QueueDeclare(
		TenantQueueName,
		true,  // durable
		false, // auto-eliminated
		false, // exclusive
		false, // no-wait
		nil,   // args
	)
	if err != nil {
		return err
	}

	// Bindear cola al exchange
	err = ch.QueueBind(
		q.Name,
		"tenant.#",
		TenantExchangeName,
		false,
		nil,
	)
	if err != nil {
		return err
	}

	logger.Info(fmt.Sprintf("Cola de tenants configurada: %s", TenantQueueName))
	return nil
}

// StartConsuming comienza a consumir mensajes
func (c *Consumer) StartConsuming(ctx context.Context) error {
	ch, err := c.conn.Channel()
	if err != nil {
		return err
	}
	defer ch.Close()

	// Configurar QoS
	err = ch.Qos(1, 0, false)
	if err != nil {
		return err
	}

	// Consumir mensajes
	msgs, err := ch.Consume(
		TenantQueueName,
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

	logger.Info(fmt.Sprintf("Comenzando a consumir mensajes de tenants desde %s", TenantQueueName))

	for {
		select {
		case <-ctx.Done():
			logger.Info("Deteniendo consumo de mensajes de tenants")
			return nil

		case msg, ok := <-msgs:
			if !ok {
				logger.Error("Canal de mensajes de tenants cerrado")
				return nil
			}

			// Procesar mensaje
			if err := c.processMessage(msg); err != nil {
				logger.Error(fmt.Sprintf("Error al procesar mensaje de tenant (routing key: %s): %v", msg.RoutingKey, err))
			}
			msg.Ack(false)
		}
	}
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

// Close cierra la conexión
func (c *Consumer) Close() error {
	if c.conn != nil {
		return c.conn.Close()
	}
	return nil
}
