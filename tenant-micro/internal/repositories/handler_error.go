package repositories

import (
	"errors"

	"github.com/maxskaink/proyecto-microservicios/tenant-micro/internal/models"
	"gorm.io/gorm"
)

func ParseDBError(err error) error {
	if err == nil {
		return nil
	}

	if errors.Is(err, gorm.ErrRecordNotFound) {
		return models.NotFoundError{Message: "Resource not found"}
	}

	if errors.Is(err, gorm.ErrDuplicatedKey) {
		return models.ConflictError{Message: "Resource already exists"}
	}

	// Si es error de restricción de clave foránea
	if errors.Is(err, gorm.ErrForeignKeyViolated) {
		return models.BadRequestError{Message: "Foreign key constraint failed"}
	}

	// Verificar si es un error de validación
	if errors.Is(err, gorm.ErrInvalidData) {
		return models.BadRequestError{Message: "Invalid data provided"}
	}

	// Para cualquier otro error, devolver un error interno del servidor
	return models.InternalServerError{Message: "Database error: " + err.Error()}
}
