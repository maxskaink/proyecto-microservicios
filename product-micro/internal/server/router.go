package server

import (
	"github.com/gin-gonic/gin"
	"github.com/maxskaink/proyecto-microservicios/product-micro/internal/controllers"
	"github.com/maxskaink/proyecto-microservicios/product-micro/internal/middleware"
)

// RegisterRoutes registra controladores en el router base.
func RegisterRoutes(r *gin.Engine) {
	// Rutas de API protegidas (con middleware de tenant)
	api := r.Group("/api")
	{
		api.Use(middleware.TenantMiddleware())
		controllers.NewProductController(ProductService).RegisterRoutes(api, AuthMiddleware)
	}
}
