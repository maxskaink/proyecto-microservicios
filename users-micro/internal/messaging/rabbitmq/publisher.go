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
func (p *Publisher) PublishMessage(ctx context.Context, routingKey string, message interface{}) error {
	ch, err := p.cm.Channel()
	if err != nil {
		return err
	}
	defer ch.Close()

	// Serializar el mensaje
	body, err := json.Marshal(message)
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
			ContentType:  "application/json",
			DeliveryMode: amqp.Persistent,
			Timestamp:    time.Now(),
			Body:         body,
		},
	)
}

// PublishUserCreated publica un evento de usuario creado
func (p *Publisher) PublishUserCreated(ctx context.Context, user interface{}) error {
	return p.PublishMessage(ctx, UserCreatedKey, user)
}

// PublishUserUpdated publica un evento de usuario actualizado
func (p *Publisher) PublishUserUpdated(ctx context.Context, user interface{}) error {
	return p.PublishMessage(ctx, UserUpdatedKey, user)
}

// PublishUserDeleted publica un evento de usuario eliminado
func (p *Publisher) PublishUserDeleted(ctx context.Context, userID string) error {
	return p.PublishMessage(ctx, UserDeletedKey, map[string]string{"id": userID})
}

// Close cierra la conexión
func (p *Publisher) Close() error {
	return p.cm.Close()
}
