package service_test

import (
	"errors"
	"github.com/google/uuid"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
	"swapp-go/cmd/internal/application/service"
	"swapp-go/cmd/internal/domain"
	"testing"
)

var username = "test_user"
var email = "test@example.com"
var password = "password123"

// Mocks
type MockUserRepository struct {
	mock.Mock
}

func (m *MockUserRepository) CreateUser(user *domain.User) error {
	return m.Called(user).Error(0)
}

func (m *MockUserRepository) UpdateUser(id uuid.UUID, fields map[string]interface{}) (*domain.User, error) {
	args := m.Called(id, fields)

	if user, ok := args.Get(0).(*domain.User); ok {
		return user, args.Error(1)
	}

	return nil, args.Error(1)
}

func (m *MockUserRepository) DeleteUser(id uuid.UUID) error {
	return m.Called(id).Error(0)
}

func (m *MockUserRepository) GetUserByID(id uuid.UUID) (*domain.User, error) {
	args := m.Called(id)

	if user, ok := args.Get(0).(*domain.User); ok {
		return user, args.Error(1)
	}

	return nil, args.Error(1)
}

func (m *MockUserRepository) GetUserByUsername(username string) (*domain.User, error) {
	args := m.Called(username)

	if user, ok := args.Get(0).(*domain.User); ok {
		return user, args.Error(1)
	}

	return nil, args.Error(1)
}

func (m *MockUserRepository) GetUserByEmail(email string) (*domain.User, error) {
	args := m.Called(email)

	if user, ok := args.Get(0).(*domain.User); ok {
		return user, args.Error(1)
	}

	return nil, args.Error(1)
}

// Helper Functions
func setupTest() (*MockUserRepository, *service.UserService) {
	mockRepo := new(MockUserRepository)
	userService := service.NewUserService(mockRepo)

	return mockRepo, userService
}

// Tests
// RegisterUser
func TestRegisterUser_Success(t *testing.T) {
	mockRepo, userService := setupTest()

	user := &domain.User{
		Username: username,
		Email:    email,
		Password: password,
	}

	mockRepo.On("GetUserByEmail", user.Email).Return(nil, errors.New("not found"))
	mockRepo.On("GetUserByUsername", user.Username).Return(nil, errors.New("not found"))
	mockRepo.On("CreateUser", mock.AnythingOfType("*domain.User")).Return(nil)

	err := userService.RegisterUser(user)

	assert.NoError(t, err)
	mockRepo.AssertExpectations(t)
}

func TestRegisterUser_EmailAlreadyExists(t *testing.T) {
	mockRepo, userService := setupTest()

	user := &domain.User{
		Username: username,
		Email:    email,
		Password: password,
	}

	mockRepo.On("GetUserByEmail", user.Email).Return(&domain.User{}, nil)

	err := userService.RegisterUser(user)

	assert.EqualError(t, err, "email already exists")
	mockRepo.AssertCalled(t, "GetUserByEmail", user.Email)
}

func TestRegisterUser_UsernameAlreadyExists(t *testing.T) {
	mockRepo, userService := setupTest()

	user := &domain.User{
		Username: username,
		Email:    email,
		Password: password,
	}

	mockRepo.On("GetUserByEmail", user.Email).Return(nil, errors.New("not found"))
	mockRepo.On("GetUserByUsername", user.Username).Return(&domain.User{}, nil)

	err := userService.RegisterUser(user)

	assert.EqualError(t, err, "username not available")
	mockRepo.AssertCalled(t, "GetUserByUsername", user.Username)
}

// UpdateUser
func TestUpdateUser_Success(t *testing.T) {
	mockRepo := new(MockUserRepository)
	userService := service.NewUserService(mockRepo)

	userID := uuid.New()
	existingUser := &domain.User{
		ID:       userID,
		Username: "old_user",
		Email:    "olduser@email.com",
	}
	updatedFields := map[string]interface{}{
		"username": "new_user",
		"email":    "newuser@email.com",
	}
	updatedUser := &domain.User{
		ID:       userID,
		Username: "new_user",
		Email:    "newuser@email.com",
	}

	mockRepo.On("GetUserByID", userID).Return(existingUser, nil)
	mockRepo.On("UpdateUser", userID, updatedFields).Return(updatedUser, nil)

	user, err := userService.UpdateUser(userID, updatedFields)
	assert.NoError(t, err)
	assert.Equal(t, updatedUser.Username, user.Username)
	assert.Equal(t, updatedUser.Email, user.Email)

	mockRepo.AssertExpectations(t)
}

func TestUpdateUser_UpdateFails(t *testing.T) {
	mockRepo := new(MockUserRepository)
	userService := service.NewUserService(mockRepo)

	id := uuid.New()
	fields := map[string]interface{}{
		"username": "wrong_user",
	}

	mockRepo.
		On("GetUserByID", id).
		Return(&domain.User{ID: id}, nil)

	mockRepo.
		On("UpdateUser", id, fields).
		Return(nil, errors.New("update error"))

	user, err := userService.UpdateUser(id, fields)

	assert.Error(t, err)
	assert.Nil(t, user)

	mockRepo.AssertExpectations(t)
}

// GetUserBy (ID|Email|Username)
func TestGetUserByID(t *testing.T) {
	mockRepo, userService := setupTest()

	userID := uuid.New()
	expectedUser := &domain.User{ID: userID, Username: username, Email: email, Password: password}

	mockRepo.On("GetUserByID", userID).Return(expectedUser, nil)

	result, err := userService.GetUserByID(userID)

	assert.NoError(t, err)
	assert.Equal(t, expectedUser, result)
	mockRepo.AssertExpectations(t)
}

func TestGetUserByEmail(t *testing.T) {
	mockRepo, userService := setupTest()

	expectedUser := &domain.User{Email: email, Username: username, Password: password}

	mockRepo.On("GetUserByEmail", email).Return(expectedUser, nil)

	result, err := userService.GetUserByEmail(email)

	assert.NoError(t, err)
	assert.Equal(t, expectedUser, result)
	mockRepo.AssertExpectations(t)
}

func TestGetUserByUsername(t *testing.T) {
	mockRepo, userService := setupTest()

	expectedUser := &domain.User{Username: username, Email: email, Password: password}

	mockRepo.On("GetUserByUsername", username).Return(expectedUser, nil)

	result, err := userService.GetUserByUsername(username)

	assert.NoError(t, err)
	assert.Equal(t, expectedUser, result)
	mockRepo.AssertExpectations(t)
}

// DeleteUser
func TestDeleteUser_Success(t *testing.T) {
	mockRepo, userService := setupTest()
	userID := uuid.New()

	mockRepo.On("DeleteUser", userID).Return(nil)

	err := userService.DeleteUser(userID)
	assert.NoError(t, err)

	mockRepo.AssertExpectations(t)
}

func TestDeleteUser_Failure(t *testing.T) {
	mockRepo, userService := setupTest()

	userID := uuid.New()
	expectedErr := errors.New("delete failed")

	mockRepo.On("DeleteUser", userID).Return(expectedErr)

	err := userService.DeleteUser(userID)

	assert.Error(t, err)
	assert.Equal(t, expectedErr, err)

	mockRepo.AssertExpectations(t)
}
