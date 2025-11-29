package controllers

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/maxskaink/proyecto-microservicios/shipping-micro/internal/domain"
	"github.com/maxskaink/proyecto-microservicios/shipping-micro/internal/dto"
	"github.com/maxskaink/proyecto-microservicios/shipping-micro/internal/middleware"
	ordersvc "github.com/maxskaink/proyecto-microservicios/shipping-micro/internal/services/order"
)

type OrderController struct {
	svc *ordersvc.Service
}

func NewOrderController(svc *ordersvc.Service) *OrderController { return &OrderController{svc: svc} }

func (oc *OrderController) Register(rg *gin.RouterGroup) {
	rg.POST("/orders", oc.createOrder)
	rg.GET("/orders/producer/:id_producer", oc.GetOrdersByProducer)
	rg.GET("/orders", oc.listUserOrders)
	rg.GET("/orders/:id", oc.getByID)
	rg.PUT("/orders/:id/status", oc.updateStatus)
}

func (oc *OrderController) createOrder(c *gin.Context) {
	tenantID := middleware.GetTenantFromContext(c)
	userID, err := middleware.GetUserIDFromContext(c)
	if err != nil {
		handleUserError(c, err)
		return
	}

	var req dto.CreateOrderRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		handleUserError(c, err)
		return
	}
	order, shipping, err := oc.svc.CreateFromCart(userID, req.ShippingAddress, tenantID)
	if err != nil {
		handleUserError(c, err)
		return
	}
	c.JSON(http.StatusCreated, gin.H{"order": order, "shipping": shipping})
}

// GetOrdersByProducer devuelve todas las órdenes que contienen al menos un producto
// cuyo Product.ProducerID coincide con producerId. Sin paginación (simple).
func (oc *OrderController) GetOrdersByProducer(c *gin.Context) {
	tenantID := middleware.GetTenantFromContext(c)
	if tenantID == "" {
		handleUserError(c, domain.BadRequestError{Message: "tenant no especificado"})
		return
	}

	producerID := c.Param("id_producer")
	if producerID == "" {
		handleUserError(c, domain.BadRequestError{Message: "producerId requerido"})
		return
	}

	orders, err := oc.svc.GetOrdersByProducer(producerID, tenantID)
	if err != nil {
		handleUserError(c, err)
		return
	}

	c.JSON(http.StatusOK, orders)
}

func (oc *OrderController) listUserOrders(c *gin.Context) {
	tenantID := middleware.GetTenantFromContext(c)
	userID, err := middleware.GetUserIDFromContext(c)

	if err != nil {
		handleUserError(c, err)
		return
	}

	statusStr := c.Query("status")

	if statusStr == "" {
		statusStr = string(domain.OrderStatusPending)
	}

	if !domain.IsValidOrderStatus(statusStr) {
		c.JSON(http.StatusBadRequest, dto.ErrorResponse{Error: "invalid status"})
		return
	}

	res, err := oc.svc.GetByUserID(userID, domain.OrderStatus(statusStr), tenantID)
	if err != nil {
		handleUserError(c, err)
		return
	}
	c.JSON(http.StatusOK, res)
}

func (oc *OrderController) getByID(c *gin.Context) {
	tenantID := middleware.GetTenantFromContext(c)

	id := c.Param("id")

	if id == "" {
		handleUserError(c, domain.BadRequestError{Message: "The id is missing"})
		return
	}
	order, err := oc.svc.GetByID(id, tenantID)

	if err != nil {
		handleUserError(c, err)
		return
	}

	c.JSON(http.StatusOK, order)
}

func (oc *OrderController) updateStatus(c *gin.Context) {
	tenantID := middleware.GetTenantFromContext(c)
	id := c.Param("id")
	if id == "" {
		c.JSON(http.StatusBadRequest, dto.ErrorResponse{Error: "invalid id"})
		return
	}

	userID, err := middleware.GetUserIDFromContext(c)
	if err != nil {
		handleUserError(c, err)
		return
	}

	var req dto.UpdateOrderStatusRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, dto.ErrorResponse{Error: "invalid body"})
		return
	}
	if err := oc.svc.UpdateStatus(id, req.Status, userID, tenantID); err != nil {
		c.JSON(http.StatusBadRequest, dto.ErrorResponse{Error: err.Error()})
		return
	}
	c.Status(http.StatusNoContent)
}
