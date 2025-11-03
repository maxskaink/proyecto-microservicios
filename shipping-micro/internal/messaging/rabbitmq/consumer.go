package rabbitmq

import (
	"context"
	"fmt"

	"github.com/maxskaink/proyecto-microservicios/shipping-micro/pkg/logger"
	amqp "github.com/rabbitmq/amqp091-go"
)

// Exchanges y colas que escuchamos en shipping
const (
	UserExchangeName = "user_events"
	UserExchangeType = "topic"
	UserQueueName    = "shipping_user_events"

	TenantExchangeName = "tenants_events"
	TenantExchangeType = "topic"
	TenantQueueName    = "shipping_tenant_events"

	ProductExchangeName = "product_events"
	ProductExchangeType = "topic"
	ProductQueueName    = "shipping_product_events"
)

type Dispatcher interface {
	ProcessMessage(data []byte, routingKey string) error
}

type Consumer struct {
	cm         *ConnectionManager
	dispatcher Dispatcher
}

func NewConsumer(cm *ConnectionManager, dispatcher Dispatcher) (*Consumer, error) {
	c := &Consumer{cm: cm, dispatcher: dispatcher}
	if err := c.setupUserQueue(); err != nil {
		return nil, err
	}
	if err := c.setupTenantQueue(); err != nil {
		return nil, err
	}
	if err := c.setupProductQueue(); err != nil {
		return nil, err
	}
	return c, nil
}

func (c *Consumer) setupUserQueue() error {
	ch, err := c.cm.Channel()
	if err != nil {
		return err
	}
	defer ch.Close()
	if err := ch.ExchangeDeclare(UserExchangeName, UserExchangeType, true, false, false, false, nil); err != nil {
		return err
	}
	q, err := ch.QueueDeclare(UserQueueName, true, false, false, false, nil)
	if err != nil {
		return err
	}
	if err := ch.QueueBind(q.Name, "user.#", UserExchangeName, false, nil); err != nil {
		return err
	}
	logger.Info(fmt.Sprintf("Cola de usuarios configurada: %s", UserQueueName))
	return nil
}

func (c *Consumer) setupTenantQueue() error {
	ch, err := c.cm.Channel()
	if err != nil {
		return err
	}
	defer ch.Close()
	if err := ch.ExchangeDeclare(TenantExchangeName, TenantExchangeType, true, false, false, false, nil); err != nil {
		return err
	}
	q, err := ch.QueueDeclare(TenantQueueName, true, false, false, false, nil)
	if err != nil {
		return err
	}
	if err := ch.QueueBind(q.Name, "tenant.#", TenantExchangeName, false, nil); err != nil {
		return err
	}
	logger.Info(fmt.Sprintf("Cola de tenants configurada: %s", TenantQueueName))
	return nil
}

func (c *Consumer) setupProductQueue() error {
	ch, err := c.cm.Channel()
	if err != nil {
		return err
	}
	defer ch.Close()
	if err := ch.ExchangeDeclare(ProductExchangeName, ProductExchangeType, true, false, false, false, nil); err != nil {
		return err
	}
	q, err := ch.QueueDeclare(ProductQueueName, true, false, false, false, nil)
	if err != nil {
		return err
	}
	if err := ch.QueueBind(q.Name, "product.#", ProductExchangeName, false, nil); err != nil {
		return err
	}
	logger.Info(fmt.Sprintf("Cola de productos configurada: %s", ProductQueueName))
	return nil
}

func (c *Consumer) StartConsuming(ctx context.Context) error {
	go func() {
		if err := c.startUserConsumer(ctx); err != nil {
			logger.Error(fmt.Sprintf("Error user consumer: %v", err))
		}
	}()
	go func() {
		if err := c.startTenantConsumer(ctx); err != nil {
			logger.Error(fmt.Sprintf("Error tenant consumer: %v", err))
		}
	}()
	go func() {
		if err := c.startProductConsumer(ctx); err != nil {
			logger.Error(fmt.Sprintf("Error product consumer: %v", err))
		}
	}()
	return nil
}

func (c *Consumer) startUserConsumer(ctx context.Context) error { return c.consume(ctx, UserQueueName) }
func (c *Consumer) startTenantConsumer(ctx context.Context) error {
	return c.consume(ctx, TenantQueueName)
}
func (c *Consumer) startProductConsumer(ctx context.Context) error {
	return c.consume(ctx, ProductQueueName)
}

func (c *Consumer) consume(ctx context.Context, queue string) error {
	ch, err := c.cm.Channel()
	if err != nil {
		return err
	}
	defer ch.Close()
	if err := ch.Qos(1, 0, false); err != nil {
		return err
	}
	msgs, err := ch.Consume(queue, "", false, false, false, false, nil)
	if err != nil {
		return err
	}
	logger.Info(fmt.Sprintf("Consumiendo mensajes desde %s", queue))
	for {
		select {
		case <-ctx.Done():
			return nil
		case msg, ok := <-msgs:
			if !ok {
				return nil
			}
			if err := c.processMessage(msg); err != nil {
				logger.Error(fmt.Sprintf("Error procesando mensaje (%s): %v", msg.RoutingKey, err))
				msg.Ack(false)
			} else {
				msg.Ack(false)
			}
		}
	}
}

func (c *Consumer) processMessage(msg amqp.Delivery) error {
	if c.dispatcher == nil {
		return fmt.Errorf("dispatcher no configurado")
	}
	return c.dispatcher.ProcessMessage(msg.Body, msg.RoutingKey)
}

func (c *Consumer) Close() error { return c.cm.Close() }
