package controllers

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/maxskaink/proyecto-microservicios/tenant-micro/internal/dto"
	"github.com/maxskaink/proyecto-microservicios/tenant-micro/internal/services"
)

type TenantController struct {
	service *services.TenantService
}

func NewTenantController(service *services.TenantService) *TenantController {
	return &TenantController{service: service}
}

func (c *TenantController) CreateTenant(ctx *gin.Context) {
	var req dto.CreateTenantRequest

	if err := ctx.ShouldBindJSON(&req); err != nil {
		ctx.JSON(http.StatusBadRequest, dto.ErrorResponse{Error: err.Error()})
		return
	}

	tenant, err := c.service.CreateTenant(ctx.Request.Context(), req)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, dto.ErrorResponse{Error: err.Error()})
		return
	}

	ctx.JSON(http.StatusCreated, tenant)
}

func (c *TenantController) GetTenant(ctx *gin.Context) {
	tenantID := ctx.Param("id")

	tenant, err := c.service.GetTenant(tenantID)
	if err != nil {
		ctx.JSON(http.StatusNotFound, dto.ErrorResponse{Error: "Tenant not found"})
		return
	}

	ctx.JSON(http.StatusOK, tenant)
}

func (c *TenantController) GetAllTenants(ctx *gin.Context) {
	tenants, err := c.service.GetAllTenants()
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, dto.ErrorResponse{Error: err.Error()})
		return
	}

	ctx.JSON(http.StatusOK, tenants)
}

func (c *TenantController) DeleteTenant(ctx *gin.Context) {
	tenantID := ctx.Param("id")

	if err := c.service.DeleteTenant(ctx.Request.Context(), tenantID); err != nil {
		ctx.JSON(http.StatusInternalServerError, dto.ErrorResponse{Error: err.Error()})
		return
	}

	ctx.JSON(http.StatusOK, gin.H{"message": "Tenant deleted"})
}
