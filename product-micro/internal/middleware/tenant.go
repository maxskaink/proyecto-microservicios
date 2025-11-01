package middleware

import (
	"github.com/gin-gonic/gin"
	"github.com/maxskaink/proyecto-microservicios/product-micro/pkg/logger"
)

const (
	// TenantHeader es el nombre del header donde se envía la información del tenant
	TenantHeader = "X-Tenant-ID"
	// TenantContextKey es la clave utilizada para almacenar el tenant en el contexto de Gin
	TenantContextKey = "tenant_id"
)

// TenantMiddleware es un middleware que extrae el tenant del header y lo almacena en el contexto
func TenantMiddleware() gin.HandlerFunc {
	return func(c *gin.Context) {
		// Obtener el tenant del header
		tenantID := c.GetHeader(TenantHeader)

		// Si no hay tenant en el header, retornar error
		if tenantID == "" {
			logger.Error("Request sin tenant identificado")
			c.AbortWithStatusJSON(400, gin.H{"error": "Tenant no especificado"})
			return
		}

		// Almacenar el tenant en el contexto
		c.Set(TenantContextKey, tenantID)
		logger.Info("Tenant identificado: " + tenantID)

		c.Next()
	}
}

// GetTenantFromContext obtiene el tenant del contexto de la solicitud
func GetTenantFromContext(c *gin.Context) string {
	tenant, exists := c.Get(TenantContextKey)
	if !exists {
		return ""
	}
	if tenantStr, ok := tenant.(string); ok {
		return tenantStr
	}
	return ""
}
