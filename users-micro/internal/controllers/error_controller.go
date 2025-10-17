package controllers

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/maxskaink/proyecto-microservicios/users-micro/internal/domain"
	"github.com/maxskaink/proyecto-microservicios/users-micro/internal/dto"
	"github.com/maxskaink/proyecto-microservicios/users-micro/pkg/logger"
)

func handleUserError(c *gin.Context, err error) {
	switch err.(type) {
	case domain.NotFoundError:
		c.JSON(dto.NewErrorDTO(http.StatusNotFound, err.Error()))
	case domain.BadRequestError:
		c.JSON(dto.NewErrorDTO(http.StatusBadRequest, err.Error()))
	case domain.UnauthorizedError:
		c.JSON(dto.NewErrorDTO(http.StatusUnauthorized, err.Error()))
	case domain.ConflictError:
		c.JSON(dto.NewErrorDTO(http.StatusConflict, err.Error()))
	default:
		logger.Error("Internal Server Error:" + err.Error())
		c.JSON(dto.NewErrorDTO(http.StatusInternalServerError, "Internal Server Error"))
	}
}
