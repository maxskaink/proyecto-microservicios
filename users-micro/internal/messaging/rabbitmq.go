package messaging

import (
	"time"

	amqp "github.com/rabbitmq/amqp091-go"
)

type ConnectionManager struct {
	conn *amqp.Connection
	url  string
}

func NewConnectionManager(url string) (*ConnectionManager, error) {
	cm := &ConnectionManager{url: url}
	if err := cm.connect(); err != nil {
		return nil, err
	}
	return cm, nil
}

func (m *ConnectionManager) connect() error {
	var err error
	// reconexión simple con backoff
	for i := 0; i < 5; i++ {
		m.conn, err = amqp.Dial(m.url)
		if err == nil {
			return nil
		}
		time.Sleep(time.Duration(i+1) * time.Second)
	}
	return err
}

func (m *ConnectionManager) Channel() (*amqp.Channel, error) {
	return m.conn.Channel()
}

func (m *ConnectionManager) Close() error {
	if m.conn != nil {
		return m.conn.Close()
	}
	return nil
}
