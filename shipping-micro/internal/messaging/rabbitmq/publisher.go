package rabbitmq

import (
	"context"
	"encoding/json"
	"time"

	amqp "github.com/rabbitmq/amqp091-go"
)

const (
	OrdersExchangeName    = "orders_events"
	OrdersExchangeType    = "topic"
	OrderCreatedKey       = "order.created"
	OrderStatusChangedKey = "order.status_changed"
	OrderPaidKey          = "order.paid"

	ShippingExchangeName = "shipping_events"
	ShippingExchangeType = "topic"
	ShippingCreatedKey   = "shipping.created"
)

type EventMessage struct {
	EventType string                 `json:"event_type"`
	Timestamp time.Time              `json:"timestamp"`
	TenantID  string                 `json:"tenant_id"`
	Data      map[string]interface{} `json:"data"`
}

type Publisher struct{ cm *ConnectionManager }

func NewPublisher(cm *ConnectionManager) (*Publisher, error) {
	p := &Publisher{cm: cm}
	if err := p.setupExchanges(); err != nil {
		return nil, err
	}
	return p, nil
}

func (p *Publisher) setupExchanges() error {
	ch, err := p.cm.Channel()
	if err != nil {
		return err
	}
	defer ch.Close()
	if err := ch.ExchangeDeclare(OrdersExchangeName, OrdersExchangeType, true, false, false, false, nil); err != nil {
		return err
	}
	if err := ch.ExchangeDeclare(ShippingExchangeName, ShippingExchangeType, true, false, false, false, nil); err != nil {
		return err
	}
	return nil
}

func (p *Publisher) publish(ctx context.Context, exchange, routingKey string, payload interface{}, tenantID string) error {
	ch, err := p.cm.Channel()
	if err != nil {
		return err
	}
	defer ch.Close()
	msg := EventMessage{EventType: routingKey, Timestamp: time.Now(), TenantID: tenantID, Data: map[string]interface{}{"payload": payload}}
	body, err := json.Marshal(msg)
	if err != nil {
		return err
	}
	ctx, cancel := context.WithTimeout(ctx, 5*time.Second)
	defer cancel()
	return ch.PublishWithContext(ctx, exchange, routingKey, false, false, amqp.Publishing{
		ContentType:  "application/json",
		DeliveryMode: amqp.Persistent,
		Timestamp:    time.Now(),
		Headers:      amqp.Table{"X-Tenant-ID": tenantID, "X-Event-Type": routingKey},
		Body:         body,
	})
}

func (p *Publisher) PublishOrderCreated(ctx context.Context, order interface{}, tenantID string) error {
	return p.publish(ctx, OrdersExchangeName, OrderCreatedKey, order, tenantID)
}

func (p *Publisher) PublishOrderStatusChanged(ctx context.Context, payload interface{}, tenantID string) error {
	return p.publish(ctx, OrdersExchangeName, OrderStatusChangedKey, payload, tenantID)
}

func (p *Publisher) PublishOrderPaid(ctx context.Context, payload interface{}, tenantID string) error {
	return p.publish(ctx, OrdersExchangeName, OrderPaidKey, payload, tenantID)
}

func (p *Publisher) PublishShippingCreated(ctx context.Context, shipping interface{}, tenantID string) error {
	return p.publish(ctx, ShippingExchangeName, ShippingCreatedKey, shipping, tenantID)
}

func (p *Publisher) Close() error { return p.cm.Close() }
