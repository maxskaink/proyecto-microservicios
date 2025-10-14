package db

import (
	"context"
	"fmt"
	"log"
	"os"

	db_models "github.com/maxskaink/proyecto-microservicios/users-micro/internal/db/models"
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
		&db_models.ProfileDB{},
	); err != nil {
		log.Fatalf("Error al migrar las tablas: %v", err)
	}

	return &GormDBProvider{db: db}, nil
}

// DB implementa el método de la interfaz DBProvider
func (p *GormDBProvider) DB(ctx context.Context) (*gorm.DB, error) {
	return p.db.WithContext(ctx), nil
}
