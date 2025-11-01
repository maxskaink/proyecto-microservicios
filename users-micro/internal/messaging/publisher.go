package messaging

import (
	"context"
)

// Publisher define la interfaz para un publicador de mensajes
type Publisher interface {
	PublishUserCreated(ctx context.Context, user interface{}, tenantID string) error
	PublishUserUpdated(ctx context.Context, user interface{}, tenantID string) error
	PublishUserDeleted(ctx context.Context, userID string, tenantID string) error
	Close() error
}
