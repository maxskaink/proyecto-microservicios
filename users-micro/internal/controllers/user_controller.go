package controllers

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/maxskaink/proyecto-microservicios/users-micro/internal/dto"
	"github.com/maxskaink/proyecto-microservicios/users-micro/internal/services"
)

// UserController maneja endpoints de usuarios.
// No implementa lógica de negocio; delega a servicios.
type UserController struct {
	UserService services.UserService
}

func NewUserController(userService services.UserService) *UserController {
	return &UserController{
		UserService: userService,
	}
}

// RegisterRoutes registra rutas HTTP relacionadas a usuarios.
func (uc *UserController) RegisterRoutes(rg *gin.RouterGroup, auth gin.HandlerFunc) {
	users := rg.Group("/users")
	{
		users.GET("/:id", auth, uc.GetUserByID)
		users.GET("/uid/:id", auth, uc.GetUserByUUID)
		users.PUT("/:id", auth, uc.UpdateUser)
		users.DELETE("/:id", auth, uc.DeleteUser)
	}
}

// GetUserByID maneja GET /users/:id.
func (uc *UserController) GetUserByID(c *gin.Context) {
	id := c.Param("id")
	user, err := uc.UserService.GetUserByID(id)
	if err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "Usuario no encontrado"})
		return
	}
	c.JSON(http.StatusOK, user)
}

func (uc *UserController) GetUserByUUID(c *gin.Context) {
	id := c.Param("id")
	user, err := uc.UserService.GetUserByUUID(id)
	if err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "Usuario no encontrado con uid " + id})
		return
	}
	c.JSON(http.StatusOK, user)
}

// UpdateUser maneja PUT /users/:id.
func (uc *UserController) UpdateUser(c *gin.Context) {
	id := c.Param("id")
	var userRequest dto.UserRequest
	if err := c.ShouldBindJSON(&userRequest); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Datos inválidos"})
		return
	}

	userResponse, err := uc.UserService.UpdateUser(id, userRequest)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Error al actualizar el usuario"})
		return
	}
	c.JSON(http.StatusOK, userResponse)
}

// DeleteUser maneja DELETE /users/:id.
func (uc *UserController) DeleteUser(c *gin.Context) {
	id := c.Param("id")
	if err := uc.UserService.DeleteUser(id); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Error al eliminar el usuario"})
		return
	}
	c.JSON(http.StatusNoContent, nil)
}
