package rabbitmq

import (
	"fmt"
	"time"

	"github.com/maxskaink/proyecto-microservicios/shipping-micro/pkg/logger"
	amqp "github.com/rabbitmq/amqp091-go"
)

type ConnectionManager struct {
	conn      *amqp.Connection
	config    *Config
	connected bool
}

func NewConnectionManager(config *Config) (*ConnectionManager, error) {
	if config == nil {
		config = DefaultConfig()
	}
	cm := &ConnectionManager{config: config}
	if err := cm.Connect(); err != nil {
		return nil, err
	}
	return cm, nil
}

func (cm *ConnectionManager) Connect() error {
	var err error
	for i := 0; i < cm.config.ReconnectRetries; i++ {
		logger.Info(fmt.Sprintf("Intentando conectar a RabbitMQ (intento %d)...", i+1))
		cm.conn, err = amqp.Dial(cm.config.URL)
		if err == nil {
			cm.connected = true
			logger.Info("Conexión exitosa a RabbitMQ")
			go cm.handleDisconnect()
			return nil
		}
		logger.Error(fmt.Sprintf("Error al conectar a RabbitMQ: %v", err))
		time.Sleep(cm.config.ReconnectInterval)
	}
	return fmt.Errorf("no se pudo establecer conexión con RabbitMQ después de %d intentos", cm.config.ReconnectRetries)
}

func (cm *ConnectionManager) handleDisconnect() {
	ch := make(chan *amqp.Error)
	cm.conn.NotifyClose(ch)
	err := <-ch
	cm.connected = false
	logger.Error(fmt.Sprintf("Conexión a RabbitMQ cerrada: %v", err))
	backoff := cm.config.ReconnectInterval
	for i := 0; i < cm.config.ReconnectRetries; i++ {
		time.Sleep(backoff)
		if err := cm.Connect(); err == nil {
			return
		}
		backoff *= 2
	}
	logger.Error("No se pudo restablecer la conexión a RabbitMQ")
}

func (cm *ConnectionManager) Channel() (*amqp.Channel, error) {
	if !cm.connected || cm.conn == nil {
		if err := cm.Connect(); err != nil {
			return nil, err
		}
	}
	return cm.conn.Channel()
}

func (cm *ConnectionManager) Close() error {
	if cm.conn != nil && cm.connected {
		return cm.conn.Close()
	}
	return nil
}

func (cm *ConnectionManager) IsConnected() bool { return cm.connected && cm.conn != nil }
