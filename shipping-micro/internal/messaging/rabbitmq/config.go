package rabbitmq

import (
	"os"
	"time"
)

type Config struct {
	URL               string
	ReconnectInterval time.Duration
	ReconnectRetries  int
}

func DefaultConfig() *Config {
	url := os.Getenv("RABBITMQ_URL")
	if url == "" {
		url = "amqp://guest:guest@localhost:5672/"
	}
	return &Config{URL: url, ReconnectInterval: 2 * time.Second, ReconnectRetries: 5}
}
