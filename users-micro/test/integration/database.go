package integration_test

import (
	"context"
	"fmt"
	"os"
	"time"

	db_models "github.com/maxskaink/proyecto-microservicios/users-micro/internal/db/models"
	tc "github.com/testcontainers/testcontainers-go"
	"github.com/testcontainers/testcontainers-go/wait"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
)

// SetupTestDB arranca un contenedor de Postgres para pruebas, ejecuta migraciones y
// devuelve la instancia *gorm.DB junto con una función de limpieza.
func SetupTestDB(ctx context.Context) (*gorm.DB, func() error, error) {
	// Opciones del contenedor (API genérica)
	req := tc.ContainerRequest{
		Image:        "postgres:15-alpine",
		Env:          map[string]string{"POSTGRES_DB": "testdb", "POSTGRES_USER": "test", "POSTGRES_PASSWORD": "secret"},
		ExposedPorts: []string{"5432/tcp"},
		WaitingFor:   wait.ForListeningPort("5432/tcp"),
	}

	// Deshabilitar el reaper (ryuk) en entornos donde no puede arrancar
	_ = os.Setenv("TESTCONTAINERS_RYUK_DISABLED", "true")

	pgC, err := tc.GenericContainer(ctx, tc.GenericContainerRequest{
		ContainerRequest: req,
		Started:          true,
	})
	if err != nil {
		return nil, nil, fmt.Errorf("no se pudo iniciar el contenedor Postgres: %w", err)
	}

	// Construir DSN
	host, err := pgC.Host(ctx)
	if err != nil {
		_ = pgC.Terminate(ctx)
		return nil, nil, err
	}
	port, err := pgC.MappedPort(ctx, "5432/tcp")
	if err != nil {
		_ = pgC.Terminate(ctx)
		return nil, nil, err
	}

	dsn := fmt.Sprintf("host=%s port=%s user=%s password=%s dbname=%s sslmode=disable", host, port.Port(), "test", "secret", "testdb")

	// Esperar un poco extra por si acaso
	time.Sleep(500 * time.Millisecond)

	db, err := gorm.Open(postgres.Open(dsn), &gorm.Config{})
	if err != nil {
		_ = pgC.Terminate(ctx)
		return nil, nil, fmt.Errorf("error al conectar con la base de datos de prueba: %w", err)
	}

	// Ejecutar migraciones necesarias para las pruebas
	if err := db.AutoMigrate(&db_models.UserDB{}, &db_models.ProfileDB{}); err != nil {
		_ = pgC.Terminate(ctx)
		return nil, nil, fmt.Errorf("error al ejecutar migraciones: %w", err)
	}

	cleanup := func() error {
		return pgC.Terminate(ctx)
	}

	return db, cleanup, nil
}
