package server

import (
	"github.com/gin-gonic/gin"
	"github.com/maxskaink/proyecto-microservicios/users-micro/internal/controllers"
)

// RegisterRoutes registra controladores en el router base.
func RegisterRoutes(r *gin.Engine) {
	api := r.Group("/api/v1")
	// TODO: agregar AuthMiddleware() cuando se integre Firebase

	controllers.NewUserController().RegisterRoutes(api)
	controllers.NewOrderController().RegisterRoutes(api)
	controllers.NewProfileController().RegisterRoutes(api)
	controllers.NewTestController().RegisterRoutes(api, FirebaseAuthMiddleware())
}
