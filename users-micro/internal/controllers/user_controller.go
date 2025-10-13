package controllers

import "github.com/gin-gonic/gin"

// UserController maneja endpoints de usuarios.
// No implementa lógica de negocio; delega a servicios.
type UserController struct {
	// userService UserService (interfaz definida en services)
}

func NewUserController() *UserController {
	return &UserController{}
}

// RegisterRoutes registra rutas HTTP relacionadas a usuarios.
func (uc *UserController) RegisterRoutes(rg *gin.RouterGroup) {
	users := rg.Group("/users")
	{
		// TODO: añadir handlers: GET /:id, POST /, PUT /:id, DELETE /:id, etc.
		_ = users
	}
}
