package server

import (
	"github.com/gin-gonic/gin"
	"github.com/maxskaink/proyecto-microservicios/users-micro/internal/controllers"
	"github.com/maxskaink/proyecto-microservicios/users-micro/internal/middleware"
)

// RegisterRoutes registra controladores en el router base.
func RegisterRoutes(r *gin.Engine) {
	// Rutas protegidas con tenant (ya tienen middleware aplicado en server.go)
	api := r.Group("api")
	api.Use(middleware.TenantMiddleware(DB))

	controllers.NewUserController(UserService).RegisterRoutes(api, AuthMiddleware)
}
