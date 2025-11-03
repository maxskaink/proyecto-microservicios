package controllers

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/maxskaink/proyecto-microservicios/shipping-micro/internal/dto"
	"github.com/maxskaink/proyecto-microservicios/shipping-micro/internal/middleware"
	shippingsvc "github.com/maxskaink/proyecto-microservicios/shipping-micro/internal/services/shipping"
)

type ShippingController struct {
	svc *shippingsvc.Service
}

func NewShippingController(svc *shippingsvc.Service) *ShippingController {
	return &ShippingController{svc: svc}
}

func (sc *ShippingController) Register(rg *gin.RouterGroup) {
	rg.GET("/shippings", sc.list)
	rg.GET("/shippings/:id", sc.getByID)
	rg.GET("/shippings/order/:order_id", sc.getByOrderID)
	rg.PUT("/shippings/:id/status", sc.updateStatus)
}

func (sc *ShippingController) list(c *gin.Context) {
	tenantID := middleware.GetTenantFromContext(c)
	res, err := sc.svc.List(tenantID)
	if err != nil {
		c.JSON(http.StatusInternalServerError, dto.ErrorResponse{Error: err.Error()})
		return
	}
	c.JSON(http.StatusOK, res)
}

func (sc *ShippingController) getByID(c *gin.Context) {
	tenantID := middleware.GetTenantFromContext(c)
	id := c.Param("id")
	if id == "" {
		c.JSON(http.StatusBadRequest, dto.ErrorResponse{Error: "invalid id"})
		return
	}
	sh, err := sc.svc.GetByID(id, tenantID)
	if err != nil {
		c.JSON(http.StatusNotFound, dto.ErrorResponse{Error: err.Error()})
		return
	}
	c.JSON(http.StatusOK, sh)
}

func (sc *ShippingController) getByOrderID(c *gin.Context) {
	tenantID := middleware.GetTenantFromContext(c)
	id := c.Param("order_id")
	if id == "" {
		c.JSON(http.StatusBadRequest, dto.ErrorResponse{Error: "invalid id"})
		return
	}
	sh, err := sc.svc.GetByOrderID(id, tenantID)
	if err != nil {
		c.JSON(http.StatusNotFound, dto.ErrorResponse{Error: err.Error()})
		return
	}
	c.JSON(http.StatusOK, sh)
}

func (sc *ShippingController) updateStatus(c *gin.Context) {
	tenantID := middleware.GetTenantFromContext(c)
	id := c.Param("id")
	if id == "" {
		c.JSON(http.StatusBadRequest, dto.ErrorResponse{Error: "invalid id"})
		return
	}
	var req dto.UpdateShippingStatusRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, dto.ErrorResponse{Error: "invalid body"})
		return
	}
	if err := sc.svc.UpdateStatus(id, req.Status, tenantID); err != nil {
		c.JSON(http.StatusBadRequest, dto.ErrorResponse{Error: err.Error()})
		return
	}
	c.Status(http.StatusNoContent)
}
