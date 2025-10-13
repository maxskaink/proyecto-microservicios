package controllers

import "github.com/gin-gonic/gin"

// OrderController maneja endpoints de pedidos (orders).
type OrderController struct{}

func NewOrderController() *OrderController { return &OrderController{} }

// RegisterRoutes registra rutas HTTP relacionadas a orders.
func (oc *OrderController) RegisterRoutes(rg *gin.RouterGroup) {
	orders := rg.Group("/orders")
	{
		// TODO: añadir handlers: GET /:id, GET /user/:userId, POST /, etc.
		_ = orders
	}
}
