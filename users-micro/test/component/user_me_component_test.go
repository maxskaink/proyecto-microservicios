package component_test

import (
	"context"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"os"
	"testing"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/stretchr/testify/require"
	tc "github.com/testcontainers/testcontainers-go"
	"github.com/testcontainers/testcontainers-go/wait"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"

	// users-micro internals
	"github.com/maxskaink/proyecto-microservicios/users-micro/internal/controllers"
	db_models "github.com/maxskaink/proyecto-microservicios/users-micro/internal/db/models"
	"github.com/maxskaink/proyecto-microservicios/users-micro/internal/db/repositories"
	"github.com/maxskaink/proyecto-microservicios/users-micro/internal/db/tenant"
	"github.com/maxskaink/proyecto-microservicios/users-micro/internal/dto"
	"github.com/maxskaink/proyecto-microservicios/users-micro/internal/middleware"
	"github.com/maxskaink/proyecto-microservicios/users-micro/internal/services"
)

// setupTestDB levanta un Postgres efímero y ejecuta migraciones mínimas
func setupTestDB(t *testing.T) (*gorm.DB, func()) {
	t.Helper()
	ctx := context.Background()

	// Deshabilitar el reaper (ryuk) en entornos donde no puede arrancar
	_ = os.Setenv("TESTCONTAINERS_RYUK_DISABLED", "true")

	req := tc.ContainerRequest{
		Image:        "postgres:15-alpine",
		Env:          map[string]string{"POSTGRES_DB": "testdb", "POSTGRES_USER": "test", "POSTGRES_PASSWORD": "secret"},
		ExposedPorts: []string{"5432/tcp"},
		WaitingFor:   wait.ForListeningPort("5432/tcp").WithStartupTimeout(60 * time.Second),
	}

	// Evitar problemas con Ryuk en algunos entornos CI (no se usa logger custom)

	pgC, err := tc.GenericContainer(ctx, tc.GenericContainerRequest{ContainerRequest: req, Started: true})
	require.NoError(t, err)

	host, err := pgC.Host(ctx)
	require.NoError(t, err)
	port, err := pgC.MappedPort(ctx, "5432/tcp")
	require.NoError(t, err)

	dsn := "host=" + host + " port=" + port.Port() + " user=test password=secret dbname=testdb sslmode=disable"
	// Pequeña espera adicional
	time.Sleep(500 * time.Millisecond)
	db, err := gorm.Open(postgres.Open(dsn), &gorm.Config{})
	require.NoError(t, err)

	// Migraciones base en public
	require.NoError(t, db.AutoMigrate(&db_models.UserDB{}, &db_models.ProfileDB{}))

	cleanup := func() {
		_ = pgC.Terminate(ctx)
	}
	return db, cleanup
}

// fakeAuthMiddleware simula la autenticación y auto-creación de usuario (como haría FirebaseAuthMiddleware)
func fakeAuthMiddleware(svc services.UserService) gin.HandlerFunc {
	return func(c *gin.Context) {
		tenantID := c.GetHeader("X-Tenant-ID")
		if tenantID == "" {
			c.AbortWithStatusJSON(http.StatusBadRequest, gin.H{"error": "Tenant no especificado"})
			return
		}

		// Datos simulados de token verificado
		uid := "uid-component-test"
		email := "component@test.local"
		name := "component"

		// Si no existe, crearlo
		if _, err := svc.GetUserByUUID(uid, tenantID); err != nil {
			_, err = svc.CreateUser(dto.UserRequest{FirebaseUID: uid, Email: email, Name: name}, tenantID)
			if err != nil {
				c.AbortWithStatusJSON(http.StatusInternalServerError, gin.H{"error": "No se pudo crear usuario de prueba"})
				return
			}
		}

		c.Set("uid", uid)
		c.Set("email", email)
		c.Set("name", name)
		c.Next()
	}
}

func Test_Component_GetMe_OK(t *testing.T) {
	db, cleanup := setupTestDB(t)
	defer cleanup()

	// Crear schema del tenant y migrar dentro del schema
	tenantDB := tenant.NewTenantDB(db)
	const tenantID = "tenant_a"
	require.NoError(t, tenantDB.CreateTenantSchema(tenantID))

	// Repo + Service
	repo := repositories.NewUserRepository(db, *tenantDB)
	svc := services.NewUserServiceWithoutPublisher(repo)

	// Router con rutas del controlador y middleware fake
	gin.SetMode(gin.TestMode)
	r := gin.New()
	api := r.Group("/api")
	api.Use(middleware.TenantMiddleware(db))
	controllers.NewUserController(svc).RegisterRoutes(api, fakeAuthMiddleware(svc))

	// Ejecutar petición
	req := httptest.NewRequest(http.MethodGet, "/api/users/me", nil)
	req.Header.Set("Authorization", "Bearer fake-token")
	req.Header.Set("X-Tenant-ID", tenantID)
	w := httptest.NewRecorder()
	r.ServeHTTP(w, req)

	require.Equal(t, http.StatusOK, w.Code)

	var body map[string]any
	require.NoError(t, json.Unmarshal(w.Body.Bytes(), &body))
	require.Equal(t, "component@test.local", body["email"])
	require.Equal(t, "uid-component-test", body["firebaseUID"])
}
