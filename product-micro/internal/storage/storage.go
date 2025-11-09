package storage

import (
	"context"
	"time"
)

// ObjectStorage define la interfaz mínima para interactuar con un almacenamiento S3-compatible.
// Implementaciones (p. ej., R2) deben cumplirla para permitir mocking y tests.
type ObjectStorage interface {
	PresignPut(ctx context.Context, bucket, key, contentType string, expires time.Duration) (uploadURL string, publicURL string, err error)
	HeadObject(ctx context.Context, bucket, key string) (ObjectMetadata, error)
	CopyObject(ctx context.Context, bucket, sourceKey, destKey string) error
	DeleteObject(ctx context.Context, bucket, key string) error
}

// ObjectMetadata es una vista mínima del objeto para validaciones.
type ObjectMetadata struct {
	ContentLength int64
	ContentType   string
	ETag          string
	LastModified  time.Time
}
