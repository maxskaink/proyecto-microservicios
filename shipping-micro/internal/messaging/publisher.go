package messaging

import "context"

// Publisher para eventos salientes del shipping micro
type Publisher interface {
	PublishOrderCreated(ctx context.Context, order interface{}, tenantID string) error
	PublishOrderStatusChanged(ctx context.Context, payload interface{}, tenantID string) error
	PublishOrderPaid(ctx context.Context, payload interface{}, tenantID string) error
	PublishShippingCreated(ctx context.Context, shipping interface{}, tenantID string) error
	Close() error
}
