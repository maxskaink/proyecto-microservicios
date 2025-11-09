package storage

import (
	"context"

	"github.com/maxskaink/proyecto-microservicios/product-micro/internal/config"
)

// NewFromConfig construye un ObjectStorage basado en la configuración (R2).
func NewFromConfig(ctx context.Context, cfg *config.Config) (ObjectStorage, error) {
	return NewR2Storage(ctx, cfg.R2Endpoint, cfg.R2Region, cfg.R2AccessKeyID, cfg.R2SecretAccessKey, cfg.R2PublicBaseURL)
}
