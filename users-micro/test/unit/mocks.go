package unit_test

import (
	"github.com/maxskaink/proyecto-microservicios/users-micro/internal/db/repositories"
	"github.com/maxskaink/proyecto-microservicios/users-micro/internal/dto"
	"github.com/stretchr/testify/mock"
)

// Crear un mock del repositorio
type MockUserRepository struct {
	mock.Mock
}

func NewMockUserService() repositories.UserRepository {
	return &MockUserRepository{}
}

func (m *MockUserRepository) Create(u *dto.UserRequest, tenantID string) (*dto.UserResponse, error) {
	args := m.Called(u, tenantID)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*dto.UserResponse), args.Error(1)
}

func (m *MockUserRepository) FindByID(id string, tenantID string) (*dto.UserResponse, error) {
	args := m.Called(id, tenantID)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*dto.UserResponse), args.Error(1)
}

func (m *MockUserRepository) FindByUUID(uuid string, tenantID string) (*dto.UserResponse, error) {
	args := m.Called(uuid, tenantID)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*dto.UserResponse), args.Error(1)
}

func (m *MockUserRepository) Update(id string, u *dto.UserRequest, tenantID string) (*dto.UserResponse, error) {
	args := m.Called(id, u, tenantID)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*dto.UserResponse), args.Error(1)
}

func (m *MockUserRepository) UpdateRol(id string, rol string, tenantID string) (*dto.UserResponse, error) {
	args := m.Called(id, rol, tenantID)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*dto.UserResponse), args.Error(1)
}

func (m *MockUserRepository) Delete(id string, tenantID string) error {
	args := m.Called(id, tenantID)
	return args.Error(0)
}
