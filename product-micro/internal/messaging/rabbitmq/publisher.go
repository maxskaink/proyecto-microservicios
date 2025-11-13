package rabbitmq

import (
	"context"
	"encoding/json"
	"fmt"
	"time"

	"github.com/maxskaink/proyecto-microservicios/product-micro/internal/dto"
	msgEvents "github.com/maxskaink/proyecto-microservicios/product-micro/internal/messaging/events"
	"github.com/maxskaink/proyecto-microservicios/product-micro/pkg/logger"
	amqp "github.com/rabbitmq/amqp091-go"
)

const (
	// Product events exchange
	ProductExchangeName = "product_events"
	ProductExchangeType = "topic"
)

// productPublisher publica eventos de productos en RabbitMQ
type productPublisher struct {
	cm *ConnectionManager
}

// NewProductPublisher crea un nuevo publisher para eventos de productos
func NewProductPublisher(config *Config) (*productPublisher, error) {
	cm, err := NewConnectionManager(config)
	if err != nil {
		return nil, err
	}
	p := &productPublisher{cm: cm}
	if err := p.ensureExchange(); err != nil {
		_ = cm.Close()
		return nil, err
	}
	return p, nil
}

// ensureExchange declara el exchange de productos si no existe
func (p *productPublisher) ensureExchange() error {
	ch, err := p.cm.Channel()
	if err != nil {
		return err
	}
	defer ch.Close()

	return ch.ExchangeDeclare(
		ProductExchangeName,
		ProductExchangeType,
		true,  // durable
		false, // auto-deleted
		false, // internal
		false, // no-wait
		nil,   // args
	)
}

// Close cierra la conexión del publisher
func (p *productPublisher) Close() error {
	return p.cm.Close()
}

// outgoingEvent es la forma de evento que publicamos (alineado con otros servicios)
type outgoingEvent struct {
	EventType msgEvents.EventType    `json:"event_type"`
	Timestamp time.Time              `json:"timestamp"`
	TenantID  string                 `json:"tenant_id"`
	Data      map[string]interface{} `json:"data"`
}

func buildProductEvent(eventType msgEvents.EventType, tenantID string, product dto.ProductDTOResponse) outgoingEvent {
	return outgoingEvent{
		EventType: eventType,
		Timestamp: time.Now().UTC(),
		TenantID:  tenantID,
		Data: map[string]interface{}{
			"payload": product,
		},
	}
}

func (p *productPublisher) publish(routingKey string, evt outgoingEvent) error {
	ch, err := p.cm.Channel()
	if err != nil {
		return err
	}
	defer ch.Close()

	body, err := json.Marshal(evt)
	if err != nil {
		return err
	}

	logger.Info(fmt.Sprintf("Publicando evento %s para tenant %s", routingKey, evt.TenantID))

	return ch.PublishWithContext(
		context.Background(),
		ProductExchangeName,
		routingKey,
		false, // mandatory
		false, // immediate
		amqp.Publishing{
			ContentType:  "application/json",
			DeliveryMode: amqp.Persistent,
			Body:         body,
			Timestamp:    evt.Timestamp,
			Type:         routingKey,
		},
	)
}

// PublishProductCreated publica un evento de producto creado
func (p *productPublisher) PublishProductCreated(product dto.ProductDTOResponse, tenantID string) error {
	evt := buildProductEvent(msgEvents.ProductCreated, tenantID, product)
	return p.publish(string(msgEvents.ProductCreated), evt)
}

// PublishProductUpdated publica un evento de producto actualizado
func (p *productPublisher) PublishProductUpdated(product dto.ProductDTOResponse, tenantID string) error {
	evt := buildProductEvent(msgEvents.ProductUpdated, tenantID, product)
	return p.publish(string(msgEvents.ProductUpdated), evt)
}

// PublishProductStockUpdated publica un evento de stock de producto actualizado
func (p *productPublisher) PublishProductStockUpdated(product dto.ProductDTOResponse, tenantID string) error {
	evt := buildProductEvent(msgEvents.ProductStockUpdated, tenantID, product)
	return p.publish(string(msgEvents.ProductStockUpdated), evt)
}
