package controllers

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/maxskaink/proyecto-microservicios/users-micro/internal/domain"
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
		users.GET("/me", auth, uc.GetInfoUser)
		users.PUT("/:id", auth, uc.UpdateUser)
		users.PATCH("/:id/rol", auth, uc.UpdateRolUser)

	}

}

// GetUserByID maneja GET /users/:id.
func (uc *UserController) GetInfoUser(c *gin.Context) {
	user_uid := c.GetString("uid") //Should have token because the middleware

	if user_uid == "" {
		handleUserError(c, domain.InternalServerError{Message: "uid vacio, no deberia de haber entrado sin uid"})
		return
	}

	user, err := uc.UserService.GetUserByUUID(user_uid)
	if err != nil {
		handleUserError(c, err)
		return
	}
	c.JSON(http.StatusOK, user)
}

// UpdateUser maneja PUT /users/:id.
func (uc *UserController) UpdateUser(c *gin.Context) {
	id := c.Param("id")
	var userRequest dto.UserRequest
	if err := c.ShouldBindJSON(&userRequest); err != nil {
		handleUserError(c, domain.BadRequestError{Message: "Los datos son invalidos: " + err.Error()})
		return
	}

	uid_requester := c.GetString("uid") // Should have it for the auth

	if uid_requester == "" {
		handleUserError(c, domain.InternalServerError{Message: "uid, vacio no deberia de haber entrado sin uid"})
		return
	}

	userResponse, err := uc.UserService.UpdateUser(id, userRequest, uid_requester)
	if err != nil {
		handleUserError(c, err)
		return
	}
	c.JSON(http.StatusOK, userResponse)
}

func (uc *UserController) UpdateRolUser(c *gin.Context) {
	// Obtener ID del usuario a modificar
	id := c.Param("id")

	// Definir la estructura para recibir los datos
	type RolUpdateRequest struct {
		Rol string `json:"rol" binding:"required"`
	}

	// Verificar que el usuario que solicita el cambio esté autenticado
	uid_requester := c.GetString("uid")
	if uid_requester == "" {
		handleUserError(c, domain.InternalServerError{Message: "uid vacío, no debería haber entrado sin uid"})
		return
	}

	// Parsear la solicitud
	var request RolUpdateRequest
	if err := c.ShouldBindJSON(&request); err != nil {
		handleUserError(c, domain.BadRequestError{Message: "Formato inválido: " + err.Error()})
		return
	}
	// Actualizar el rol
	userResponse, err := uc.UserService.UpdateRol(id, domain.UserRole(request.Rol), uid_requester)
	if err != nil {
		handleUserError(c, err)
		return
	}

	c.JSON(http.StatusOK, userResponse)
}
