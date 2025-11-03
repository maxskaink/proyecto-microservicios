package unit_test

import (
	"errors"
	"testing"

	"github.com/maxskaink/proyecto-microservicios/users-micro/internal/domain"
	"github.com/maxskaink/proyecto-microservicios/users-micro/internal/dto"
	"github.com/maxskaink/proyecto-microservicios/users-micro/internal/services"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
)

var tenantID = "tenantA"

func TestCreateUser(t *testing.T) {
	// Arrange
	mockRepo := new(MockUserRepository)
	service := services.NewUserServiceWithoutPublisher(mockRepo)

	testUserRequest := dto.UserRequest{
		Email:       "test@example.com",
		Name:        "Test User",
		FirebaseUID: "firebase123", 
	}

	expectedResponse := &dto.UserResponse{
		ID:    "1",
		Email: "test@example.com",
		Name:  "Test User",
	}

	mockRepo.On("Create", mock.Anything).Return(expectedResponse, nil)

	// Act 
	response, err := service.CreateUser(testUserRequest, tenantID)

	// Assert
	assert.NoError(t, err)
	assert.Equal(t, expectedResponse.ID, response.ID)
	assert.Equal(t, expectedResponse.Email, response.Email)
	assert.Equal(t, expectedResponse.Name, response.Name)
	mockRepo.AssertExpectations(t)
}

func TestCreateUserValidationError(t *testing.T) {
	// Arrange
	mockRepo := new(MockUserRepository)
	service := services.NewUserServiceWithoutPublisher(mockRepo)

	// Caso de prueba: email vacío
	testUserRequest := dto.UserRequest{
		Email:       "",
		Name:        "Test User",
		FirebaseUID: "firebase123",
	}

	// Act
	_, err := service.CreateUser(testUserRequest, tenantID)

	// Assert
	assert.Error(t, err)
	var badRequestError domain.BadRequestError
	isBadRequest := errors.As(err, &badRequestError)
	assert.True(t, isBadRequest)
	mockRepo.AssertNotCalled(t, "Create")
}

func TestGetUserByID(t *testing.T) {
	// Arrange
	mockRepo := new(MockUserRepository)
	service := services.NewUserServiceWithoutPublisher(mockRepo)

	userID := "1"
	expectedUser := &dto.UserResponse{
		ID:    userID,
		Email: "test@example.com",
		Name:  "Test User",
	}

	mockRepo.On("FindByID", userID).Return(expectedUser, nil)

	// Act
	user, err := service.GetUserByID(userID, tenantID)

	// Assert
	assert.NoError(t, err)
	assert.Equal(t, expectedUser.ID, user.ID)
	assert.Equal(t, expectedUser.Email, user.Email)
	mockRepo.AssertExpectations(t)
}

func TestGetUserByIDNotFound(t *testing.T) {
	// Arrange
	mockRepo := new(MockUserRepository)
	service := services.NewUserServiceWithoutPublisher(mockRepo)

	userID := "999"
	mockRepo.On("FindByID", userID).Return(nil, errors.New("usuario no encontrado"))

	// Act
	_, err := service.GetUserByID(userID, tenantID)

	// Assert
	assert.Error(t, err)
	mockRepo.AssertExpectations(t)
}

func TestUpdateUser(t *testing.T) {
	// Arrange
	mockRepo := new(MockUserRepository)
	service := services.NewUserServiceWithoutPublisher(mockRepo)

	userID := "1"
	firebaseUID := "firebase123"

	// El usuario encontrado en la BD
	foundUser := &dto.UserResponse{
		ID:          userID,
		Email:       "old@example.com",
		Name:        "Old Name",
		FirebaseUID: firebaseUID,
	}

	// La actualización que queremos hacer
	updateRequest := dto.UserRequest{
		Email:       "new@example.com",
		Name:        "New Name",
		FirebaseUID: firebaseUID,
	}

	// El resultado esperado después de la actualización
	expectedResponse := &dto.UserResponse{
		ID:          userID,
		Email:       "new@example.com",
		Name:        "New Name",
		FirebaseUID: firebaseUID,
	}

	// Configurar el comportamiento del mock
	mockRepo.On("FindByID", userID).Return(foundUser, nil)
	mockRepo.On("Update", userID, &updateRequest).Return(expectedResponse, nil)

	// Act
	response, err := service.UpdateUser(userID, updateRequest, firebaseUID, tenantID)

	// Assert
	assert.NoError(t, err)
	assert.Equal(t, expectedResponse.Email, response.Email)
	assert.Equal(t, expectedResponse.Name, response.Name)
	mockRepo.AssertExpectations(t)
}

func TestUpdateUserNoPermission(t *testing.T) {
	// Arrange
	mockRepo := new(MockUserRepository)
	service := services.NewUserServiceWithoutPublisher(mockRepo)

	userID := "1"
	userFirebaseUID := "firebase123"
	differentFirebaseUID := "firebase456"

	foundUser := &dto.UserResponse{
		ID:          userID,
		Email:       "test@example.com",
		Name:        "Test User",
		FirebaseUID: userFirebaseUID,
	}

	updateRequest := dto.UserRequest{
		Email: "new@example.com",
		Name:  "New Name",
	}

	mockRepo.On("FindByID", userID).Return(foundUser, nil)

	// Act
	_, err := service.UpdateUser(userID, updateRequest, differentFirebaseUID, tenantID)

	// Assert
	assert.Error(t, err)
	assert.Contains(t, err.Error(), "no tiene permisos")
	mockRepo.AssertNotCalled(t, "Update")
}

func TestDeleteUser(t *testing.T) {
	// Arrange
	mockRepo := new(MockUserRepository)
	service := services.NewUserServiceWithoutPublisher(mockRepo)

	userID := "1"
	mockRepo.On("Delete", userID).Return(nil)

	// Act
	err := service.DeleteUser(userID, tenantID)

	// Assert
	assert.NoError(t, err)
	mockRepo.AssertExpectations(t)
}
