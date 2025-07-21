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

type MockUserRepository struct {
	mock.Mock
}

func (m *MockUserRepository) CreateUser(user *domain.User) error {
	args := m.Called(user)

	return args.Error(0)
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

func setupTest() (*MockUserRepository, *service.UserService) {
	mockRepo := new(MockUserRepository)
	userService := service.NewUserService(mockRepo)

	return mockRepo, userService
}

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
