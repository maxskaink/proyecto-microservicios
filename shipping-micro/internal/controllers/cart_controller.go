package controllers

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/maxskaink/proyecto-microservicios/shipping-micro/internal/domain"
	"github.com/maxskaink/proyecto-microservicios/shipping-micro/internal/dto"
	"github.com/maxskaink/proyecto-microservicios/shipping-micro/internal/middleware"
	cartsvc "github.com/maxskaink/proyecto-microservicios/shipping-micro/internal/services/cart"
)

type CartController struct {
	svc *cartsvc.Service
}

func NewCartController(svc *cartsvc.Service) *CartController { return &CartController{svc: svc} }

func (cc *CartController) Register(rg *gin.RouterGroup) {
	rg.POST("/cart/items", cc.addItem)
	rg.GET("/cart", cc.getCart)
	rg.PUT("/cart/items/:id", cc.updateItem)
	rg.DELETE("/cart/items/:id", cc.deleteItem)
	rg.DELETE("/cart", cc.clearCart)
}

func (cc *CartController) addItem(c *gin.Context) {
	userID, err := middleware.GetUserIDFromContext(c)
	if err != nil {
		handleUserError(c, err)
		return
	}
	tenantID := middleware.GetTenantFromContext(c)

	var req dto.AddToCartRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		handleUserError(c, err)
		return
	}

	item, err := cc.svc.AddItem(userID, req.ProductID, req.Quantity, tenantID)
	if err != nil {
		handleUserError(c, err)
		return
	}
	c.JSON(http.StatusCreated, item)
}

func (cc *CartController) getCart(c *gin.Context) {
	userID, err := middleware.GetUserIDFromContext(c)
	if err != nil {
		handleUserError(c, err)
		return
	}
	tenantID := middleware.GetTenantFromContext(c)

	items, err := cc.svc.GetUserCart(userID, tenantID)
	if err != nil {
		handleUserError(c, err)
		return
	}
	c.JSON(http.StatusOK, items)
}

func (cc *CartController) updateItem(c *gin.Context) {
	tenantID := middleware.GetTenantFromContext(c)
	itemID := c.Param("id")
	if itemID == "" {
		handleUserError(c, domain.BadRequestError{Message: "item id is required"})
		return
	}
	var req dto.UpdateCartItemRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		handleUserError(c, err)
		return
	}
	if err := cc.svc.UpdateItem(itemID, req.Quantity, tenantID); err != nil {
		handleUserError(c, err)
		return
	}
	c.Status(http.StatusNoContent)
}

func (cc *CartController) deleteItem(c *gin.Context) {
	tenantID := middleware.GetTenantFromContext(c)
	itemID := c.Param("id")
	if itemID == "" {
		handleUserError(c, domain.BadRequestError{Message: "item id is required"})
		return
	}
	if err := cc.svc.DeleteItem(itemID, tenantID); err != nil {
		handleUserError(c, err)
		return
	}
	c.Status(http.StatusNoContent)
}

func (cc *CartController) clearCart(c *gin.Context) {
	userID, err := middleware.GetUserIDFromContext(c)
	if err != nil {
		handleUserError(c, err)
		return
	}
	tenantID := middleware.GetTenantFromContext(c)
	if err := cc.svc.ClearCart(userID, tenantID); err != nil {
		handleUserError(c, err)
		return
	}
	c.Status(http.StatusNoContent)
}
