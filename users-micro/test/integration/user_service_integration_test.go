package integration_test

import (
	"context"
	"testing"

	"github.com/google/uuid"
	"github.com/maxskaink/proyecto-microservicios/users-micro/internal/db/repositories"
	"github.com/maxskaink/proyecto-microservicios/users-micro/internal/dto"
	"github.com/maxskaink/proyecto-microservicios/users-micro/internal/services"
	"github.com/stretchr/testify/require"
)

func setupService(t *testing.T) (services.UserService, func()) {
	ctx := context.Background()
	db, cleanup, err := SetupTestDB(ctx)
	require.NoError(t, err)

	repo := repositories.NewUserRepository(db)
	svc := services.NewUserService(repo)

	return svc, func() { _ = cleanup() }
}

func TestService_CreateUser_Success(t *testing.T) {
	svc, cleanup := setupService(t)
	defer cleanup()

	req := dto.UserRequest{
		FirebaseUID: "uid-svc-create",
		Email:       "svc.create@example.com",
		Name:        "Service Create",
	}

	created, err := svc.CreateUser(req)
	require.NoError(t, err)
	require.NotEmpty(t, created.ID)
	require.Equal(t, req.FirebaseUID, created.FirebaseUID)
	require.Equal(t, req.Email, created.Email)
	require.Equal(t, req.Name, created.Name)
}

func TestService_GetUserByID_Success(t *testing.T) {
	svc, cleanup := setupService(t)
	defer cleanup()

	req := dto.UserRequest{
		FirebaseUID: "uid-svc-get",
		Email:       "svc.get@example.com",
		Name:        "Service Get",
	}

	created, err := svc.CreateUser(req)
	require.NoError(t, err)

	got, err := svc.GetUserByID(created.ID)
	require.NoError(t, err)
	require.Equal(t, created.ID, got.ID)
	require.Equal(t, created.Email, got.Email)
}

func TestService_GetUserByID_NotFound(t *testing.T) {
	svc, cleanup := setupService(t)
	defer cleanup()

	nonexistent := uuid.NewString()
	_, err := svc.GetUserByID(nonexistent)
	require.Error(t, err)
}

func TestService_UpdateUser_Success(t *testing.T) {
	svc, cleanup := setupService(t)
	defer cleanup()

	req := dto.UserRequest{
		FirebaseUID: "uid-svc-update",
		Email:       "svc.update@example.com",
		Name:        "Service ToUpdate",
	}

	created, err := svc.CreateUser(req)
	require.NoError(t, err)

	// Update name
	req.Name = "Service Updated"
	updated, err := svc.UpdateUser(created.ID, req, req.FirebaseUID)
	require.NoError(t, err)
	require.Equal(t, "Service Updated", updated.Name)
	require.Equal(t, created.ID, updated.ID)
}

func TestService_UpdateUser_Unauthorized(t *testing.T) {
	svc, cleanup := setupService(t)
	defer cleanup()

	req := dto.UserRequest{
		FirebaseUID: "uid-svc-auth",
		Email:       "svc.auth@example.com",
		Name:        "Service Auth",
	}

	created, err := svc.CreateUser(req)
	require.NoError(t, err)

	// Attempt update with a different firebaseUID (should be unauthorized)
	req.Name = "Service Hacked"
	_, err = svc.UpdateUser(created.ID, req, "different-uid")
	require.Error(t, err)
}
