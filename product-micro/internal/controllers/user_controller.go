package controllers

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/maxskaink/proyecto-microservicios/users-micro/internal/domain"
)

// ProductController maneja endpoints de productos.
// Punto de acceso para los endpoinst realacionados a los productos
type ProductController struct {
	//Deberia de tener los servicios que manejan los casos de uso
}

func NewUserController() *ProductController {
	return &ProductController{}
}

// RegisterRoutes registra rutas HTTP relacionadas a usuarios.
func (uc *ProductController) RegisterRoutes(rg *gin.RouterGroup, auth gin.HandlerFunc) {
	users := rg.Group("/users")
	{
		users.GET("/ping", auth, uc.Pong)
	}

}

// GetUserByID maneja GET /users/:id.
func (uc *ProductController) Pong(c *gin.Context) {
	user_uid := c.GetString("uid") //Should have token because the middleware

	if user_uid == "" {
		handleUserError(c, domain.InternalServerError{Message: "uid vacio, no deberia de haber entrado sin uid"})
		return
	}

	c.String(http.StatusOK, "Pong")
}
