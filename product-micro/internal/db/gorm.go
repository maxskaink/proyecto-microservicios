package db

import (
	"context"
	"errors"
	"fmt"
	"log"
	"os"

	db_models "github.com/maxskaink/proyecto-microservicios/product-micro/internal/db/models"
	"github.com/maxskaink/proyecto-microservicios/product-micro/internal/domain"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
)

type GormDBProvider struct {
	db *gorm.DB
}

// NewGormDBProvider crea una nueva instancia de GormDBProvider
func NewGormDBProvider() (*GormDBProvider, error) {
	// Obtener la URL de conexión desde las variables de entorno
	dsn := os.Getenv("DATABASE_URL")
	if dsn == "" {
		return nil, fmt.Errorf("DATABASE_URL no está configurada")
	}

	// Conectar a la base de datos
	db, err := gorm.Open(postgres.Open(dsn), &gorm.Config{})
	if err != nil {
		return nil, fmt.Errorf("error al conectar a la base de datos: %w", err)
	}

	if err := db.AutoMigrate(
		&db_models.UserDB{},
		&db_models.ProductDB{},
	); err != nil {
		log.Fatalf("Error al migrar las tablas: %v", err)
	}

	return &GormDBProvider{db: db}, nil
}

// DB implementa el método de la interfaz DBProvider
func (p *GormDBProvider) DB(ctx context.Context) (*gorm.DB, error) {
	return p.db.WithContext(ctx), nil
}

func ParseDBError(err error) error {
	if err == nil {
		return nil
	}

	if errors.Is(err, gorm.ErrRecordNotFound) {
		return domain.NotFoundError{Message: "Resource not found"}
	}

	if errors.Is(err, gorm.ErrDuplicatedKey) {
		return domain.ConflictError{Message: "Resource already exists"}
	}

	// Si es error de restricción de clave foránea
	if errors.Is(err, gorm.ErrForeignKeyViolated) {
		return domain.BadRequestError{Message: "Foreign key constraint failed"}
	}

	// Verificar si es un error de validación
	if errors.Is(err, gorm.ErrInvalidData) {
		return domain.BadRequestError{Message: "Invalid data provided"}
	}

	// Para cualquier otro error, devolver un error interno del servidor
	return domain.InternalServerError{Message: "Database error: " + err.Error()}
}
