package server

import (
	"github.com/gin-gonic/gin"
	"github.com/maxskaink/proyecto-microservicios/users-micro/internal/controllers"
)

// RegisterRoutes registra controladores en el router base.
func RegisterRoutes(r *gin.Engine) {
	api := r.Group("")

	controllers.NewUserController(UserService).RegisterRoutes(api, AuthMiddleware)
}
