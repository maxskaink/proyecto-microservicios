package controllers

import (
	"fmt"
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/maxskaink/proyecto-microservicios/tenant-micro/internal/dto"
	"github.com/maxskaink/proyecto-microservicios/tenant-micro/internal/models"
)

func handleUserError(c *gin.Context, err error) {
	switch err.(type) {
	case models.NotFoundError:
		c.JSON(dto.NewErrorDTO(http.StatusNotFound, err.Error()))
	case models.BadRequestError:
		c.JSON(dto.NewErrorDTO(http.StatusBadRequest, err.Error()))
	case models.UnauthorizedError:
		c.JSON(dto.NewErrorDTO(http.StatusUnauthorized, err.Error()))
	case models.ConflictError:
		c.JSON(dto.NewErrorDTO(http.StatusConflict, err.Error()))
	case models.InvalidInputError:
		c.JSON(dto.NewErrorDTO(http.StatusBadRequest, err.Error()))
	default:
		fmt.Println("Internal Server Error:" + err.Error())
		c.JSON(dto.NewErrorDTO(http.StatusInternalServerError, "Internal Server Error"))
	}
}
