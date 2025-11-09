package storage

import (
	"context"
	"errors"
	"fmt"
	"net/url"
	"os"
	"time"

	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/config"
	"github.com/aws/aws-sdk-go-v2/credentials"
	"github.com/aws/aws-sdk-go-v2/service/s3"
)

// R2Storage implementa ObjectStorage para Cloudflare R2 (S3-compatible)
type R2Storage struct {
	client *s3.Client
	// publicBaseURL se usa para formar URLs públicas (cdn/dom)
	publicBaseURL string
}

// NewR2Storage crea un cliente R2. Requiere endpoint y credenciales S3 compatibles.
func NewR2Storage(ctx context.Context, endpoint, region, accessKey, secretKey, publicBaseURL string) (*R2Storage, error) {
	if endpoint == "" {
		return nil, errors.New("endpoint R2 requerido")
	}
	creds := credentials.NewStaticCredentialsProvider(accessKey, secretKey, "")
	cfg, err := config.LoadDefaultConfig(ctx,
		config.WithRegion(region),
		config.WithCredentialsProvider(creds),
	)
	if err != nil {
		return nil, fmt.Errorf("error cargando config AWS: %w", err)
	}

	resolver := aws.EndpointResolverFunc(func(service, r string) (aws.Endpoint, error) {
		return aws.Endpoint{URL: endpoint, SigningRegion: region, HostnameImmutable: true}, nil
	})
	cfg.EndpointResolver = resolver

	client := s3.NewFromConfig(cfg, func(o *s3.Options) {
		o.UsePathStyle = true // R2 requiere path style normalmente
	})

	return &R2Storage{client: client, publicBaseURL: publicBaseURL}, nil
}

// PresignPut genera una URL prefirmada para subir un objeto vía PUT.
func (r *R2Storage) PresignPut(ctx context.Context, bucket, key, contentType string, expires time.Duration) (string, string, error) {
	presigner := s3.NewPresignClient(r.client)
	in := &s3.PutObjectInput{
		Bucket:      aws.String(bucket),
		Key:         aws.String(key),
		ContentType: aws.String(contentType),
	}
	req, err := presigner.PresignPutObject(ctx, in, func(opts *s3.PresignOptions) {
		opts.Expires = expires
	})
	if err != nil {
		return "", "", fmt.Errorf("error presign put: %w", err)
	}

	publicURL := r.publicURL(bucket, key)
	return req.URL, publicURL, nil
}

// HeadObject obtiene metadatos mínimos del objeto.
func (r *R2Storage) HeadObject(ctx context.Context, bucket, key string) (ObjectMetadata, error) {
	out, err := r.client.HeadObject(ctx, &s3.HeadObjectInput{Bucket: aws.String(bucket), Key: aws.String(key)})
	if err != nil {
		return ObjectMetadata{}, fmt.Errorf("error head object: %w", err)
	}
	meta := ObjectMetadata{}
	if out.ContentLength != nil {
		meta.ContentLength = *out.ContentLength
	}
	if out.ContentType != nil {
		meta.ContentType = *out.ContentType
	}
	if out.ETag != nil {
		meta.ETag = *out.ETag
	}
	if out.LastModified != nil {
		meta.LastModified = *out.LastModified
	}
	return meta, nil
}

// CopyObject copia un objeto dentro del mismo bucket.
func (r *R2Storage) CopyObject(ctx context.Context, bucket, sourceKey, destKey string) error {
	src := fmt.Sprintf("%s/%s", bucket, sourceKey)
	_, err := r.client.CopyObject(ctx, &s3.CopyObjectInput{
		Bucket:     aws.String(bucket),
		CopySource: aws.String(src),
		Key:        aws.String(destKey),
	})
	if err != nil {
		return fmt.Errorf("error copy object: %w", err)
	}
	return nil
}

// DeleteObject elimina un objeto.
func (r *R2Storage) DeleteObject(ctx context.Context, bucket, key string) error {
	_, err := r.client.DeleteObject(ctx, &s3.DeleteObjectInput{Bucket: aws.String(bucket), Key: aws.String(key)})
	if err != nil {
		return fmt.Errorf("error delete object: %w", err)
	}
	return nil
}

// publicURL construye URL pública si se configuró R2_PUBLIC_BASE_URL; si no, genera una URL directa (no recomendada en prod).
func (r *R2Storage) publicURL(bucket, key string) string {
	if r.publicBaseURL != "" {
		return fmt.Sprintf("%s/%s", r.publicBaseURL, url.PathEscape(key))
	}
	// fallback: usar endpoint directo (requiere que el bucket sea público por políticas Cloudflare)
	endpoint := os.Getenv("R2_ENDPOINT")
	return fmt.Sprintf("%s/%s/%s", endpoint, bucket, url.PathEscape(key))
}
