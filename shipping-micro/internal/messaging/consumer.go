package messaging

import "context"

// Consumer define la interfaz para un consumidor de mensajes
// Compatible con otros micros del repo
type Consumer interface {
	StartConsuming(ctx context.Context) error
	Close() error
}
