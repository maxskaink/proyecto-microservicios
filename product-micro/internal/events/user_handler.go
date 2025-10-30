package events

import (
	"context"
	"encoding/json"
	"fmt"

	"github.com/maxskaink/proyecto-microservicios/product-micro/internal/messaging"
	"github.com/maxskaink/proyecto-microservicios/product-micro/internal/messaging/rabbitmq"
	"github.com/maxskaink/proyecto-microservicios/product-micro/pkg/logger"
)

// UserHandler maneja los eventos de usuarios
type UserHandler struct {
	consumer messaging.Consumer
}

// NewUserHandler crea un nuevo manejador de eventos de usuarios
func NewUserHandler() (*UserHandler, error) {
	factory := messaging.NewFactory(nil)

	// Crear un consumidor para eventos de usuarios
	consumer, err := factory.CreateConsumer(rabbitmq.QueueName, nil) // el handler se asignará después
	if err != nil {
		return nil, fmt.Errorf("error al crear consumidor de eventos: %w", err)
	}

	handler := &UserHandler{
		consumer: consumer,
	}

	// Asignar el handler después de crear el objeto
	consumer, err = factory.CreateConsumer(rabbitmq.QueueName, handler.handleUserEvent)
	if err != nil {
		return nil, fmt.Errorf("error al crear consumidor de eventos: %w", err)
	}
	handler.consumer = consumer

	return handler, nil
}

// Start inicia el consumo de mensajes
func (h *UserHandler) Start(ctx context.Context) error {
	return h.consumer.StartConsuming(ctx)
}

// Close cierra el consumidor
func (h *UserHandler) Close() error {
	return h.consumer.Close()
}

// handleUserEvent maneja un evento de usuario
func (h *UserHandler) handleUserEvent(data []byte) error {
	// Ejemplo de estructura para un usuario
	type User struct {
		ID          string `json:"id"`
		FirebaseUID string `json:"firebase_uid"`
		Email       string `json:"email"`
		Name        string `json:"name"`
		Rol         string `json:"rol"`
	}

	// Decodificar el mensaje
	var user User
	if err := json.Unmarshal(data, &user); err != nil {
		logger.Error(fmt.Sprintf("Error al decodificar evento de usuario: %v", err))
		return err
	}

	// Simplemente imprimir la información por ahora
	logger.Info(fmt.Sprintf("Usuario recibido: ID=%s, Email=%s, Nombre=%s, Rol=%s",
		user.ID, user.Email, user.Name, user.Rol))

	// Aquí podrías realizar acciones basadas en el evento (por ejemplo, actualizar caché, notificar, etc.)

	return nil
}
