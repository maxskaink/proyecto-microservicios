package messaging

import (
	"context"
)

// Publisher define la interfaz para un publicador de mensajes
type Publisher interface {
	PublishUserCreated(ctx context.Context, user interface{}) error
	PublishUserUpdated(ctx context.Context, user interface{}) error
	PublishUserDeleted(ctx context.Context, userID string) error
	Close() error
}
