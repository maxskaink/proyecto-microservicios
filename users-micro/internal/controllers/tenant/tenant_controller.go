package tenant

import (
	"fmt"
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/maxskaink/proyecto-microservicios/users-micro/internal/services/tenant"
	"github.com/maxskaink/proyecto-microservicios/users-micro/pkg/logger"
)

// TenantController maneja los endpoints relacionados con tenants
type TenantController struct {
	tenantService *tenant.TenantService
}

// NewTenantController crea un nuevo controlador de tenants
func NewTenantController(ts *tenant.TenantService) *TenantController {
	return &TenantController{
		tenantService: ts,
	}
}

// CreateTenantRequest estructura para crear un tenant
type CreateTenantRequest struct {
	TenantID   string `json:"tenant_id" binding:"required"`
	TenantName string `json:"tenant_name" binding:"required"`
}

// CreateTenant crea un nuevo tenant
// @Summary Crear un nuevo tenant
// @Description Crea un nuevo tenant con su respectivo schema en base de datos
// @Tags tenants
// @Accept json
// @Produce json
// @Param request body CreateTenantRequest true "Datos del tenant"
// @Success 201 {object} gin.H
// @Failure 400 {object} gin.H
// @Failure 500 {object} gin.H
// @Router /admin/tenants [post]
func (tc *TenantController) CreateTenant(c *gin.Context) {
	var req CreateTenantRequest

	// Validar el JSON de entrada
	if err := c.BindJSON(&req); err != nil {
		logger.Error(fmt.Sprintf("Error al parsear request de creación de tenant: %v", err))
		c.JSON(http.StatusBadRequest, gin.H{
			"error": "Datos de entrada inválidos",
		})
		return
	}

	// Validar que el tenant ID sea válido
	if !isValidTenantID(req.TenantID) {
		c.JSON(http.StatusBadRequest, gin.H{
			"error": "Tenant ID inválido. Solo se permiten letras, números, guiones y guiones bajos",
		})
		return
	}

	// Crear el tenant
	if err := tc.tenantService.CreateTenant(c, req.TenantID, req.TenantName); err != nil {
		logger.Error(fmt.Sprintf("Error al crear tenant: %v", err))
		c.JSON(http.StatusInternalServerError, gin.H{
			"error": "Error al crear tenant",
		})
		return
	}

	c.JSON(http.StatusCreated, gin.H{
		"message":     "Tenant creado exitosamente",
		"tenant_id":   req.TenantID,
		"tenant_name": req.TenantName,
	})
}

// DeleteTenant elimina un tenant
// @Summary Eliminar un tenant
// @Description Elimina un tenant y su schema de base de datos
// @Tags tenants
// @Param tenant_id path string true "ID del tenant"
// @Success 200 {object} gin.H
// @Failure 400 {object} gin.H
// @Failure 500 {object} gin.H
// @Router /admin/tenants/{tenant_id} [delete]
func (tc *TenantController) DeleteTenant(c *gin.Context) {
	tenantID := c.Param("tenant_id")

	// Validar que el tenant ID no esté vacío
	if tenantID == "" {
		c.JSON(http.StatusBadRequest, gin.H{
			"error": "Tenant ID es requerido",
		})
		return
	}

	// Eliminar el tenant
	if err := tc.tenantService.DeleteTenant(c, tenantID); err != nil {
		logger.Error(fmt.Sprintf("Error al eliminar tenant: %v", err))
		c.JSON(http.StatusInternalServerError, gin.H{
			"error": "Error al eliminar tenant",
		})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"message":   "Tenant eliminado exitosamente",
		"tenant_id": tenantID,
	})
}

// isValidTenantID valida que el tenant ID sea válido
// Solo permite letras, números, guiones y guiones bajos
func isValidTenantID(tenantID string) bool {
	if len(tenantID) == 0 || len(tenantID) > 64 {
		return false
	}
	for _, ch := range tenantID {
		if !((ch >= 'a' && ch <= 'z') || (ch >= 'A' && ch <= 'Z') ||
			(ch >= '0' && ch <= '9') || ch == '-' || ch == '_') {
			return false
		}
	}
	return true
}
