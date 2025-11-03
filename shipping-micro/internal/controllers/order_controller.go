package controllers

import (
	"net/http"

	"github.com/gin-gonic/gin"
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
	rg.GET("/orders", oc.listUserOrders)
	rg.GET("/orders/:id", oc.getByID)
	rg.PUT("/orders/:id/status", oc.updateStatus)
}

func (oc *OrderController) createOrder(c *gin.Context) {
	tenantID := middleware.GetTenantFromContext(c)
	userID, err := middleware.GetUserIDFromContext(c)
	if err != nil {
		c.JSON(http.StatusUnauthorized, dto.ErrorResponse{Error: "unauthorized"})
		return
	}

	var req dto.CreateOrderRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, dto.ErrorResponse{Error: "invalid body"})
		return
	}
	order, shipping, err := oc.svc.CreateFromCart(userID, req.ShippingAddress, tenantID)
	if err != nil {
		c.JSON(http.StatusBadRequest, dto.ErrorResponse{Error: err.Error()})
		return
	}
	c.JSON(http.StatusCreated, gin.H{"order": order, "shipping": shipping})
}

func (oc *OrderController) listUserOrders(c *gin.Context) {
	tenantID := middleware.GetTenantFromContext(c)
	userID, err := middleware.GetUserIDFromContext(c)
	if err != nil {
		c.JSON(http.StatusUnauthorized, dto.ErrorResponse{Error: "unauthorized"})
		return
	}

	res, err := oc.svc.GetByUserID(userID, tenantID)
	if err != nil {
		c.JSON(http.StatusInternalServerError, dto.ErrorResponse{Error: err.Error()})
		return
	}
	c.JSON(http.StatusOK, res)
}

func (oc *OrderController) getByID(c *gin.Context) {
	tenantID := middleware.GetTenantFromContext(c)
	_ = tenantID // reserved for future policies

	id := c.Param("id")
	if id == "" {
		c.JSON(http.StatusBadRequest, dto.ErrorResponse{Error: "invalid id"})
		return
	}
	order, err := oc.svc.GetByID(id, middleware.GetTenantFromContext(c))
	if err != nil {
		c.JSON(http.StatusNotFound, dto.ErrorResponse{Error: err.Error()})
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
	var req dto.UpdateOrderStatusRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, dto.ErrorResponse{Error: "invalid body"})
		return
	}
	if err := oc.svc.UpdateStatus(id, req.Status, tenantID); err != nil {
		c.JSON(http.StatusBadRequest, dto.ErrorResponse{Error: err.Error()})
		return
	}
	c.Status(http.StatusNoContent)
}
