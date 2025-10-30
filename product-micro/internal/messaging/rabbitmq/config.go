package rabbitmq

import (
	"os"
	"time"
)

// Config contiene la configuración para conectarse a RabbitMQ
type Config struct {
	URL               string
	ReconnectInterval time.Duration
	ReconnectRetries  int
}

// DefaultConfig retorna una configuración con valores predeterminados
func DefaultConfig() *Config {
	url := os.Getenv("RABBITMQ_URL")
	if url == "" {
		url = "amqp://guest:guest@localhost:5672/"
	}

	return &Config{
		URL:               url,
		ReconnectInterval: 2 * time.Second,
		ReconnectRetries:  5,
	}
}
