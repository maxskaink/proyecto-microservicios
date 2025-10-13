package controllers

import (
	"github.com/gin-gonic/gin"
)

// TestController manage some endpoints for test
type TestController struct{}

func NewTestController() *TestController { return &TestController{} }

func (tc *TestController) RegisterRoutes(rg *gin.RouterGroup, auth gin.HandlerFunc) {
	api := rg.Group("")

	if nil == auth {
		panic("Error in injection of dependencies")
	}

	api.GET("/public", tc.PublicController)
	api.GET("/private", auth, tc.PrivateController)
}

// PrivateController godoc
// @Summary Only accesible with a token
// @Security BearerAuth
// @Description  Say something when you have a token
// @Produce      json
// @Success      200  {object}  map[string]string  "ok"
// @Router       /api/v1/private [get]
func (tc *TestController) PrivateController(rg *gin.Context) {
	rg.JSON(200, gin.H{
		"message": "Everythin is fine here",
	})
}

// PublicController godoc
// @Summary is a public controller
// @Description  Say somethingc
// @Produce      json
// @Success      200  {object}  map[string]string  "ok"
// @Router       /api/v1/private [get]
func (tc *TestController) PublicController(rg *gin.Context) {
	rg.JSON(200, gin.H{
		"message": "Everythin is fine here",
	})
}
