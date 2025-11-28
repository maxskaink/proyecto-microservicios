package controllers

import (
	"fmt"
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/maxskaink/proyecto-microservicios/shipping-micro/internal/db/repositories"
	"github.com/maxskaink/proyecto-microservicios/shipping-micro/internal/domain"
	"github.com/maxskaink/proyecto-microservicios/shipping-micro/internal/dto"
	"github.com/maxskaink/proyecto-microservicios/shipping-micro/internal/middleware"
	shippingsvc "github.com/maxskaink/proyecto-microservicios/shipping-micro/internal/services/shipping"
)

type ShippingController struct {
	svc            *shippingsvc.Service
	userRepository repositories.IUserRepository
}

func NewShippingController(svc *shippingsvc.Service, userRepository repositories.IUserRepository) *ShippingController {
	return &ShippingController{svc: svc, userRepository: userRepository}
}

func (sc *ShippingController) Register(rg *gin.RouterGroup) {
	rg.GET("/shippings", sc.list)
	rg.GET("/shippings/:id", sc.getByID)
	rg.GET("/shippings/order/:order_id", sc.getByOrderID)
	rg.PUT("/shippings/:id/status", sc.updateStatus)
}

func (sc *ShippingController) list(c *gin.Context) {
	tenantID := middleware.GetTenantFromContext(c)

	//only the admin
	requester_id, err := middleware.GetUserIDFromContext(c)

	if err != nil {
		handleUserError(c, err)
		return
	}

	user, err := sc.userRepository.GetByID(requester_id, tenantID)
	if err != nil {
		handleUserError(c, err)
		return
	}
	fmt.Println(user)

	if user.Rol != "admin" {
		handleUserError(c, domain.UnauthorizedError{Message: "only admins can access"})
		return
	}

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
		handleUserError(c, domain.BadRequestError{Message: "Invalid Id"})
		return
	}
	sh, err := sc.svc.GetByID(id, tenantID)
	if err != nil {
		handleUserError(c, err)
		return
	}
	c.JSON(http.StatusOK, sh)
}

func (sc *ShippingController) getByOrderID(c *gin.Context) {
	tenantID := middleware.GetTenantFromContext(c)
	id := c.Param("order_id")
	if id == "" {
		handleUserError(c, domain.BadRequestError{Message: "order_id is needed"})
		return
	}
	sh, err := sc.svc.GetByOrderID(id, tenantID)
	if err != nil {
		handleUserError(c, err)
		return
	}
	c.JSON(http.StatusOK, sh)
}

func (sc *ShippingController) updateStatus(c *gin.Context) {
	tenantID := middleware.GetTenantFromContext(c)
	id := c.Param("id")
	if id == "" {
		handleUserError(c, domain.BadRequestError{Message: "invalid id"})
		return
	}

	//Validate if is admin the requester
	requester_id, err := middleware.GetUserIDFromContext(c)

	if err != nil {
		handleUserError(c, err)
		return
	}

	user, err := sc.userRepository.GetByID(requester_id, tenantID)
	if err != nil {
		handleUserError(c, err)
		return
	}

	if user.Rol != "admin" {
		handleUserError(c, domain.UnauthorizedError{Message: "only admins can access"})
		return
	}

	var req dto.UpdateShippingStatusRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		handleUserError(c, domain.BadRequestError{Message: "invalid body"})
		return
	}
	if err := sc.svc.UpdateStatus(id, req.Status, tenantID); err != nil {
		handleUserError(c, err)
		return
	}
	c.Status(http.StatusNoContent)
}
