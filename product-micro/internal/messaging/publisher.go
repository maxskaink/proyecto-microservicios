package messaging

import "github.com/maxskaink/proyecto-microservicios/product-micro/internal/dto"

// Publisher define la interfaz para publicar eventos de productos
type Publisher interface {
	PublishProductCreated(product dto.ProductDTOResponse, tenantID string) error
	PublishProductUpdated(product dto.ProductDTOResponse, tenantID string) error
	PublishProductStockUpdated(product dto.ProductDTOResponse, tenantID string) error
	Close() error
}
