package services_test

import (
	"errors"
	"github.com/google/uuid"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
	"swapp-go/cmd/internal/application/services"
	"swapp-go/cmd/internal/application/services/mocks"
	"swapp-go/cmd/internal/domain"
	"testing"
)

var (
	username        = "test_user"
	email           = "test@example.com"
	password        = "password123"
	phone           = "+447712345678"
	address         = "1, Main Street"
	updatedUsername = "updated_user"
	updatedEmail    = "updated@example.com"
	updatedPhone    = "+44778654321"
	updatedAddress  = "2, Main Street"
)

// Test Helpers
func setupTest() (*mocks.MockUserRepository, *services.UserService) {
	mockRepo := new(mocks.MockUserRepository)
	userService := services.NewUserService(mockRepo)

	return mockRepo, userService
}

// Tests
func TestRegisterUser_Success(t *testing.T) {
	mockRepo, userService := setupTest()

	user := &domain.User{
		Username: username,
		Email:    email,
		Phone:    &phone,
		Address:  &address,
		Password: password,
	}

	mockRepo.On("FindByEmail", user.Email).Return(nil, errors.New("not found"))
	mockRepo.On("FindByUsername", user.Username).Return(nil, errors.New("not found"))
	mockRepo.On("Create", mock.AnythingOfType("*domain.User")).Return(nil)

	err := userService.RegisterUser(user)

	assert.NoError(t, err)
	mockRepo.AssertExpectations(t)
}

func TestRegisterUser_EmailAlreadyExists(t *testing.T) {
	mockRepo, userService := setupTest()

	user := &domain.User{
		Username: username,
		Email:    email,
		Phone:    &phone,
		Address:  &address,
		Password: password,
	}

	mockRepo.On("FindByEmail", user.Email).Return(&domain.User{}, nil)

	err := userService.RegisterUser(user)

	assert.EqualError(t, err, "email already exists")
	mockRepo.AssertCalled(t, "FindByEmail", user.Email)
}

func TestRegisterUser_UsernameAlreadyExists(t *testing.T) {
	mockRepo, userService := setupTest()

	user := &domain.User{
		Username: username,
		Email:    email,
		Phone:    &phone,
		Address:  &address,
		Password: password,
	}

	mockRepo.On("FindByEmail", user.Email).Return(nil, errors.New("not found"))
	mockRepo.On("FindByUsername", user.Username).Return(&domain.User{}, nil)

	err := userService.RegisterUser(user)

	assert.EqualError(t, err, "username not available")
	mockRepo.AssertCalled(t, "FindByUsername", user.Username)
}

func TestUpdateUser_Success(t *testing.T) {
	mockRepo, userService := setupTest()

	userID := uuid.New()
	existingUser := &domain.User{
		ID:       userID,
		Username: username,
		Email:    email,
		Phone:    &phone,
		Address:  &address,
	}
	updatedFields := map[string]interface{}{
		"username": updatedUsername,
		"email":    updatedEmail,
		"phone":    updatedPhone,
		"address":  updatedAddress,
	}
	updatedUser := &domain.User{
		ID:       userID,
		Username: updatedUsername,
		Email:    updatedEmail,
		Phone:    &updatedPhone,
		Address:  &updatedAddress,
	}

	mockRepo.On("FindByID", userID).Return(existingUser, nil)
	mockRepo.On("Update", userID, updatedFields).Return(updatedUser, nil)

	user, err := userService.Update(userID, updatedFields)
	assert.NoError(t, err)
	assert.Equal(t, updatedUser.Username, user.Username)
	assert.Equal(t, updatedUser.Email, user.Email)
	assert.Equal(t, *updatedUser.Phone, *user.Phone)
	assert.Equal(t, *updatedUser.Address, *user.Address)

	mockRepo.AssertExpectations(t)
}

func TestUpdateUser_UpdateFails(t *testing.T) {
	mockRepo, userService := setupTest()

	id := uuid.New()
	fields := map[string]interface{}{
		"username": "wrong_user",
	}

	mockRepo.
		On("FindByID", id).
		Return(&domain.User{ID: id}, nil)

	mockRepo.
		On("Update", id, fields).
		Return(nil, errors.New("update error"))

	user, err := userService.Update(id, fields)

	assert.Error(t, err)
	assert.Nil(t, user)

	mockRepo.AssertExpectations(t)
}

func TestGetUserByID(t *testing.T) {
	mockRepo, userService := setupTest()

	userID := uuid.New()
	expectedUser := &domain.User{
		ID:       userID,
		Username: username,
		Email:    email,
		Phone:    &phone,
		Address:  &address,
		Password: password,
	}

	mockRepo.On("FindByID", userID).Return(expectedUser, nil)

	result, err := userService.FindByID(userID)

	assert.NoError(t, err)
	assert.Equal(t, expectedUser, result)
	mockRepo.AssertExpectations(t)
}

func TestGetUserByEmail(t *testing.T) {
	mockRepo, userService := setupTest()

	expectedUser := &domain.User{
		Username: username,
		Email:    email,
		Phone:    &phone,
		Address:  &address,
		Password: password}

	mockRepo.On("FindByEmail", email).Return(expectedUser, nil)

	result, err := userService.FindByEmail(email)

	assert.NoError(t, err)
	assert.Equal(t, expectedUser, result)
	mockRepo.AssertExpectations(t)
}

func TestGetUserByUsername(t *testing.T) {
	mockRepo, userService := setupTest()

	expectedUser := &domain.User{
		Username: username,
		Email:    email,
		Phone:    &phone,
		Address:  &address,
		Password: password,
	}

	mockRepo.On("FindByUsername", username).Return(expectedUser, nil)

	result, err := userService.FindByUsername(username)

	assert.NoError(t, err)
	assert.Equal(t, expectedUser, result)
	mockRepo.AssertExpectations(t)
}

func TestDeleteUser_Success(t *testing.T) {
	mockRepo, userService := setupTest()
	userID := uuid.New()

	mockRepo.On("Delete", userID).Return(nil)

	err := userService.Delete(userID)
	assert.NoError(t, err)

	mockRepo.AssertExpectations(t)
}

func TestDeleteUser_Failure(t *testing.T) {
	mockRepo, userService := setupTest()

	userID := uuid.New()
	expectedErr := errors.New("delete failed")

	mockRepo.On("Delete", userID).Return(expectedErr)

	err := userService.Delete(userID)

	assert.Error(t, err)
	assert.Equal(t, expectedErr, err)

	mockRepo.AssertExpectations(t)
}
