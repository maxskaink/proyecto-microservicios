package rabbitmq

import (
	"context"
	"fmt"

	"github.com/maxskaink/proyecto-microservicios/product-micro/pkg/logger"
	amqp "github.com/rabbitmq/amqp091-go"
)

// User Events
const (
	UserExchangeName = "user_events"
	UserExchangeType = "topic"
	UserCreatedKey   = "user.created"
	UserUpdatedKey   = "user.updated"
	UserDeletedKey   = "user.deleted"
	UserQueueName    = "product_user_events"
)

// Tenant Events
const (
	TenantExchangeName = "tenants_events"
	TenantExchangeType = "topic"
	TenantCreatedKey   = "tenant.created"
	TenantDeletedKey   = "tenant.deleted"
	TenantQueueName    = "product_tenant_events"
)

// Order Events
const (
	OrderExchangeName = "orders_events"
	OrderExchangeType = "topic"
	OrderPaidKey      = "order.paid"
	OrderQueueName    = "product_order_events"
)

type Consumer struct {
	cm         *ConnectionManager
	dispatcher Dispatcher
}

// Dispatcher define la interfaz para procesar mensajes
type Dispatcher interface {
	ProcessMessage(data []byte, routingKey string) error
}

// NewConsumer crea un nuevo consumidor
func NewConsumer(cm *ConnectionManager, dispatcher Dispatcher) (*Consumer, error) {
	c := &Consumer{
		cm:         cm,
		dispatcher: dispatcher,
	}

	// Configurar ambos exchanges y colas
	if err := c.setupUserQueue(); err != nil {
		return nil, err
	}

	if err := c.setupTenantQueue(); err != nil {
		return nil, err
	}

	if err := c.setupOrderQueue(); err != nil {
		return nil, err
	}

	return c, nil
}

// setupUserQueue inicializa el exchange y la cola para eventos de usuarios
func (c *Consumer) setupUserQueue() error {
	ch, err := c.cm.Channel()
	if err != nil {
		return err
	}
	defer ch.Close()

	err = ch.ExchangeDeclare(
		UserExchangeName,
		UserExchangeType,
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
		UserQueueName,
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
		UserExchangeName,
		false,
		nil,
	)

	if err != nil {
		return err
	}

	logger.Info(fmt.Sprintf("Cola de usuarios configurada: %s", UserQueueName))
	return nil
}

// setupTenantQueue inicializa el exchange y la cola para eventos de tenants
func (c *Consumer) setupTenantQueue() error {
	ch, err := c.cm.Channel()
	if err != nil {
		return err
	}
	defer ch.Close()

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

// setupOrderQueue inicializa el exchange y la cola para eventos de órdenes
func (c *Consumer) setupOrderQueue() error {
	ch, err := c.cm.Channel()
	if err != nil {
		return err
	}
	defer ch.Close()

	err = ch.ExchangeDeclare(
		OrderExchangeName,
		OrderExchangeType,
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
		OrderQueueName,
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
		"order.#",
		OrderExchangeName,
		false,
		nil,
	)

	if err != nil {
		return err
	}

	logger.Info(fmt.Sprintf("Cola de órdenes configurada: %s", OrderQueueName))
	return nil
}

// StartConsuming comienza a consumir mensajes de todos los exchanges
func (c *Consumer) StartConsuming(ctx context.Context) error {
	// Consumidor de eventos de usuario
	go func() {
		if err := c.startUserConsumer(ctx); err != nil {
			logger.Error(fmt.Sprintf("Error en consumidor de usuarios: %v", err))
		}
	}()

	// Consumidor de eventos de tenant
	go func() {
		if err := c.startTenantConsumer(ctx); err != nil {
			logger.Error(fmt.Sprintf("Error en consumidor de tenants: %v", err))
		}
	}()

	// Consumidor de eventos de órdenes
	go func() {
		if err := c.startOrderConsumer(ctx); err != nil {
			logger.Error(fmt.Sprintf("Error en consumidor de órdenes: %v", err))
		}
	}()

	return nil
}

// startUserConsumer consume mensajes de user_events
func (c *Consumer) startUserConsumer(ctx context.Context) error {
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
		UserQueueName,
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

	logger.Info(fmt.Sprintf("Comenzando a consumir mensajes de usuarios desde %s", UserQueueName))

	for {
		select {
		case <-ctx.Done():
			logger.Info("Deteniendo consumo de mensajes de usuarios")
			return nil

		case msg, ok := <-msgs:
			if !ok {
				logger.Error("Canal de mensajes de usuarios cerrado")
				return nil
			}

			if err := c.processMessage(msg); err != nil {
				logger.Error(fmt.Sprintf("Error al procesar mensaje de usuario (routing key: %s): %v", msg.RoutingKey, err))
				msg.Ack(false)
			} else {
				msg.Ack(false)
			}
		}
	}
}

// startTenantConsumer consume mensajes de tenants_events
func (c *Consumer) startTenantConsumer(ctx context.Context) error {
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

			if err := c.processMessage(msg); err != nil {
				logger.Error(fmt.Sprintf("Error al procesar mensaje de tenant (routing key: %s): %v", msg.RoutingKey, err))
				msg.Ack(false)
			} else {
				msg.Ack(false)
			}
		}
	}
}

// startOrderConsumer consume mensajes de orders_events
func (c *Consumer) startOrderConsumer(ctx context.Context) error {
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
		OrderQueueName,
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

	logger.Info(fmt.Sprintf("Comenzando a consumir mensajes de órdenes desde %s", OrderQueueName))

	for {
		select {
		case <-ctx.Done():
			logger.Info("Deteniendo consumo de mensajes de órdenes")
			return nil

		case msg, ok := <-msgs:
			if !ok {
				logger.Error("Canal de mensajes de órdenes cerrado")
				return nil
			}

			if err := c.processMessage(msg); err != nil {
				logger.Error(fmt.Sprintf("Error al procesar mensaje de orden (routing key: %s): %v", msg.RoutingKey, err))
				msg.Ack(false)
			} else {
				msg.Ack(false)
			}
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

// Close cierra el consumidor
func (c *Consumer) Close() error {
	return c.cm.Close()
}
