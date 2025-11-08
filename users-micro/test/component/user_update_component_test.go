package component_test

import (
	"bytes"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/gin-gonic/gin"
	"github.com/stretchr/testify/require"

	// users-micro internals
	"github.com/maxskaink/proyecto-microservicios/users-micro/internal/controllers"
	"github.com/maxskaink/proyecto-microservicios/users-micro/internal/db/repositories"
	"github.com/maxskaink/proyecto-microservicios/users-micro/internal/db/tenant"
	"github.com/maxskaink/proyecto-microservicios/users-micro/internal/dto"
	"github.com/maxskaink/proyecto-microservicios/users-micro/internal/middleware"
	"github.com/maxskaink/proyecto-microservicios/users-micro/internal/services"
)

// Nota: reutilizamos setupTestDB y fakeAuthMiddleware definidos en user_me_component_test.go

func Test_Component_UpdateUser_OK(t *testing.T) {
	db, cleanup := setupTestDB(t)
	defer cleanup()

	// Crear schema del tenant y migrar dentro del schema
	tenantDB := tenant.NewTenantDB(db)
	const tenantID = "tenant_a"
	require.NoError(t, tenantDB.CreateTenantSchema(tenantID))

	// Repo + Service
	repo := repositories.NewUserRepository(db, *tenantDB)
	svc := services.NewUserServiceWithoutPublisher(repo)

	// Crear usuario inicial (con el mismo uid que inyecta el fakeAuth)
	created, err := svc.CreateUser(dto.UserRequest{
		FirebaseUID: "uid-component-test",
		Email:       "component@test.local",
		Name:        "component",
	}, tenantID)
	require.NoError(t, err)

	// Router con rutas del controlador y middleware fake
	gin.SetMode(gin.TestMode)
	r := gin.New()
	api := r.Group("/api")
	api.Use(middleware.TenantMiddleware(db))
	controllers.NewUserController(svc).RegisterRoutes(api, fakeAuthMiddleware(svc))

	// Preparar payload de actualización
	bodyReq := dto.UserRequest{
		Email: "component.updated@test.local",
		Name:  "component-updated",
		Profile: &dto.ProfileRequest{
			Address:   "New Street 123",
			Phone:     "555-1234",
			AvatarURL: "new-avatar.png",
		},
	}
	payload, _ := json.Marshal(bodyReq)

	// Ejecutar petición PUT /api/users/:id
	req := httptest.NewRequest(http.MethodPut, "/api/users/"+created.ID, bytes.NewReader(payload))
	req.Header.Set("Authorization", "Bearer fake-token")
	req.Header.Set("X-Tenant-ID", tenantID)
	req.Header.Set("Content-Type", "application/json")
	w := httptest.NewRecorder()
	r.ServeHTTP(w, req)

	require.Equal(t, http.StatusOK, w.Code)

	var resp dto.UserResponse
	require.NoError(t, json.Unmarshal(w.Body.Bytes(), &resp))

	// Aserciones principales
	require.Equal(t, created.ID, resp.ID)
	require.Equal(t, "uid-component-test", resp.FirebaseUID)
	require.Equal(t, "component.updated@test.local", resp.Email)
	require.Equal(t, "component-updated", resp.Name)
	require.Equal(t, "New Street 123", resp.Profile.Address)
	require.Equal(t, "555-1234", resp.Profile.Phone)
	require.Equal(t, "new-avatar.png", resp.Profile.AvatarURL)
}
