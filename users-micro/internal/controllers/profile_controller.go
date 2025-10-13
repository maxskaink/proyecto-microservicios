package controllers

import "github.com/gin-gonic/gin"

// ProfileController maneja endpoints de perfil de usuario.
type ProfileController struct{}

func NewProfileController() *ProfileController { return &ProfileController{} }

// RegisterRoutes registra rutas HTTP relacionadas a perfiles.
func (pc *ProfileController) RegisterRoutes(rg *gin.RouterGroup) {
	profiles := rg.Group("/profiles")
	{
		// TODO: handlers: GET /:userId, POST /, PUT /:userId
		_ = profiles
	}
}
