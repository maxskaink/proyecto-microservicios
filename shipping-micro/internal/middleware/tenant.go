package middleware

import (
"fmt"

"github.com/gin-gonic/gin"
"github.com/maxskaink/proyecto-microservicios/shipping-micro/pkg/logger"
"gorm.io/gorm"
)

const (
TenantHeader     = "X-Tenant-Id"
TenantContextKey = "tenant_id"
)

// TenantMiddleware valida el tenant del header
func TenantMiddleware(db *gorm.DB) gin.HandlerFunc {
	return func(c *gin.Context) {
		tenantID := c.GetHeader("X-Tenant-ID")
		if tenantID == "" {
			tenantID = c.GetHeader("X-Tenant-Id")
		}

		if tenantID == "" {
			logger.Error("Request sin tenant identificado")
			c.AbortWithStatusJSON(400, gin.H{"error": "Tenant no especificado"})
			return
		}

		var exists bool
		query := `SELECT EXISTS(SELECT 1 FROM information_schema.schemata WHERE schema_name = $1)`
		if err := db.Raw(query, tenantID).Scan(&exists).Error; err != nil {
			logger.Error(fmt.Sprintf("Error verificando existencia del tenant '%s': %v", tenantID, err))
			c.AbortWithStatusJSON(500, gin.H{"error": "Error interno al validar tenant"})
			return
		}

		if !exists {
			logger.Error(fmt.Sprintf("Tenant '%s' no existe", tenantID))
			c.AbortWithStatusJSON(404, gin.H{"error": fmt.Sprintf("Tenant '%s' no encontrado", tenantID)})
			return
		}

		c.Set(TenantContextKey, tenantID)
		logger.Info("Tenant identificado y validado: " + tenantID)

		c.Next()
	}
}

// GetTenantFromContext obtiene el tenant del contexto
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
