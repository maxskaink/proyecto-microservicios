package rabbitmq

import (
	"context"
	"encoding/json"
	"time"

	amqp "github.com/rabbitmq/amqp091-go"
)

const (
	// Constantes para la configuración de RabbitMQ
	ExchangeName   = "user_events"
	ExchangeType   = "topic"
	UserCreatedKey = "user.created"
	UserUpdatedKey = "user.updated"
	UserDeletedKey = "user.deleted"
)

// EventMessage es la estructura base para todos los eventos
type EventMessage struct {
	EventType string                 `json:"event_type"`
	Timestamp time.Time              `json:"timestamp"`
	TenantID  string                 `json:"tenant_id"`
	Data      map[string]interface{} `json:"data"`
}

// Publisher maneja la publicación de mensajes
type Publisher struct {
	cm *ConnectionManager
}

// NewPublisher crea un nuevo publicador
func NewPublisher(cm *ConnectionManager) (*Publisher, error) {
	p := &Publisher{
		cm: cm,
	}

	// Inicializar el exchange
	if err := p.setupExchange(); err != nil {
		return nil, err
	}

	return p, nil
}

// setupExchange inicializa el exchange para los eventos de usuarios
func (p *Publisher) setupExchange() error {
	ch, err := p.cm.Channel()
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

	return err
}

// PublishMessage publica un mensaje en el exchange
func (p *Publisher) PublishMessage(ctx context.Context, routingKey string, message interface{}, tenantID string) error {
	ch, err := p.cm.Channel()
	if err != nil {
		return err
	}
	defer ch.Close()

	// Crear el mensaje con información de tenant
	eventMsg := EventMessage{
		EventType: routingKey,
		Timestamp: time.Now(),
		TenantID:  tenantID,
		Data: map[string]interface{}{
			"payload": message,
		},
	}

	// Serializar el mensaje
	body, err := json.Marshal(eventMsg)
	if err != nil {
		return err
	}

	// Publicar el mensaje
	ctx, cancel := context.WithTimeout(ctx, 5*time.Second)
	defer cancel()

	return ch.PublishWithContext(
		ctx,
		ExchangeName, // exchange
		routingKey,   // routing key
		false,        // mandatory
		false,        // immediate
		amqp.Publishing{
			ContentType: "application/json",
			DeliveryMode: amqp.Persistent,
			Timestamp:   time.Now(),
			Headers: amqp.Table{
				"X-Tenant-ID": tenantID,
				"X-Event-Type": routingKey,
			},
			Body: body,
		},
	)
}

// PublishUserCreated publica un evento de usuario creado
func (p *Publisher) PublishUserCreated(ctx context.Context, user interface{}, tenantID string) error {
	return p.PublishMessage(ctx, UserCreatedKey, user, tenantID)
}

// PublishUserUpdated publica un evento de usuario actualizado
func (p *Publisher) PublishUserUpdated(ctx context.Context, user interface{}, tenantID string) error {
	return p.PublishMessage(ctx, UserUpdatedKey, user, tenantID)
}

// PublishUserDeleted publica un evento de usuario eliminado
func (p *Publisher) PublishUserDeleted(ctx context.Context, userID string, tenantID string) error {
	return p.PublishMessage(ctx, UserDeletedKey, map[string]string{"id": userID}, tenantID)
}

// Close cierra la conexión
func (p *Publisher) Close() error {
	return p.cm.Close()
}
