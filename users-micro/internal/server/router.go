package server

import (
	"github.com/gin-gonic/gin"
	"github.com/maxskaink/proyecto-microservicios/users-micro/internal/controllers"
	tenantControllers "github.com/maxskaink/proyecto-microservicios/users-micro/internal/controllers/tenant"
)

// RegisterRoutes registra controladores en el router base.
func RegisterRoutes(r *gin.Engine) {
	// Rutas de administración (sin middleware de tenant)
	adminRoutes := r.Group("/admin")
	{
		tenantCtrl := tenantControllers.NewTenantController(TenantService)
		adminRoutes.POST("/tenants", tenantCtrl.CreateTenant)
		adminRoutes.DELETE("/tenants/:tenant_id", tenantCtrl.DeleteTenant)
	}

	// Rutas protegidas con tenant (ya tienen middleware aplicado en server.go)
	api := r.Group("/api")

	controllers.NewUserController(UserService).RegisterRoutes(api, AuthMiddleware)
}
