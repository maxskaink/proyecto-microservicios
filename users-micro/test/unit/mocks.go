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

func (m *MockUserRepository) Create(u *dto.UserRequest) (*dto.UserResponse, error) {
	args := m.Called(u)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*dto.UserResponse), args.Error(1)
}

func (m *MockUserRepository) FindByID(id string) (*dto.UserResponse, error) {
	args := m.Called(id)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*dto.UserResponse), args.Error(1)
}

func (m *MockUserRepository) FindByUUID(uuid string) (*dto.UserResponse, error) {
	args := m.Called(uuid)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*dto.UserResponse), args.Error(1)
}

func (m *MockUserRepository) Update(id string, u *dto.UserRequest) (*dto.UserResponse, error) {
	args := m.Called(id, u)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*dto.UserResponse), args.Error(1)
}

func (m *MockUserRepository) UpdateRol(id string, rol string) (*dto.UserResponse, error) {
	args := m.Called(id, rol)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*dto.UserResponse), args.Error(1)
}

func (m *MockUserRepository) Delete(id string) error {
	args := m.Called(id)
	return args.Error(0)
}
