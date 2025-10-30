package rabbitmq

import (
	"fmt"
	"time"

	"github.com/maxskaink/proyecto-microservicios/users-micro/pkg/logger"
	amqp "github.com/rabbitmq/amqp091-go"
)

// ConnectionManager maneja la conexión a RabbitMQ
type ConnectionManager struct {
	conn      *amqp.Connection
	config    *Config
	connected bool
}

// NewConnectionManager crea una nueva instancia de ConnectionManager
func NewConnectionManager(config *Config) (*ConnectionManager, error) {
	if config == nil {
		config = DefaultConfig()
	}

	cm := &ConnectionManager{
		config: config,
	}

	if err := cm.Connect(); err != nil {
		return nil, err
	}

	return cm, nil
}

// Connect establece una conexión con RabbitMQ
func (cm *ConnectionManager) Connect() error {
	var err error
	for i := 0; i < cm.config.ReconnectRetries; i++ {
		logger.Info(fmt.Sprintf("Intentando conectar a RabbitMQ (intento %d)...", i+1))
		cm.conn, err = amqp.Dial(cm.config.URL)
		if err == nil {
			cm.connected = true
			logger.Info("Conexión exitosa a RabbitMQ")
			// Configurar notificación de cierre de conexión para reintentar
			go cm.handleDisconnect()
			return nil
		}
		logger.Error(fmt.Sprintf("Error al conectar a RabbitMQ: %v", err))
		time.Sleep(cm.config.ReconnectInterval)
	}
	return fmt.Errorf("no se pudo establecer conexión con RabbitMQ después de %d intentos", cm.config.ReconnectRetries)
}

// handleDisconnect maneja la reconexión automática cuando se pierde la conexión
func (cm *ConnectionManager) handleDisconnect() {
	connClosed := make(chan *amqp.Error)
	cm.conn.NotifyClose(connClosed)

	err := <-connClosed
	cm.connected = false
	logger.Error(fmt.Sprintf("Conexión a RabbitMQ cerrada: %v", err))

	// Reintentar conexión
	backoff := cm.config.ReconnectInterval
	for i := 0; i < cm.config.ReconnectRetries; i++ {
		time.Sleep(backoff)
		if err := cm.Connect(); err == nil {
			return
		}
		backoff *= 2 // Backoff exponencial
	}

	logger.Error("No se pudo restablecer la conexión a RabbitMQ. Se requiere intervención manual.")
}

// Channel crea un nuevo canal en la conexión actual
func (cm *ConnectionManager) Channel() (*amqp.Channel, error) {
	if !cm.connected || cm.conn == nil {
		if err := cm.Connect(); err != nil {
			return nil, err
		}
	}
	return cm.conn.Channel()
}

// Close cierra la conexión
func (cm *ConnectionManager) Close() error {
	if cm.conn != nil && cm.connected {
		return cm.conn.Close()
	}
	return nil
}

// IsConnected devuelve el estado de la conexión
func (cm *ConnectionManager) IsConnected() bool {
	return cm.connected && cm.conn != nil
}
