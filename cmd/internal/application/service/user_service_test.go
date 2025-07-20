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

func TestRegisterUser_Success(t *testing.T) {
	mockRepo := new(MockUserRepository)
	userService := service.NewUserService(mockRepo)

	user := &domain.User{
		Username: "testuser",
		Email:    "test@example.com",
		Password: "password123",
	}

	mockRepo.On("GetUserByEmail", user.Email).Return(nil, errors.New("not found"))
	mockRepo.On("GetUserByUsername", user.Username).Return(nil, errors.New("not found"))
	mockRepo.On("CreateUser", mock.AnythingOfType("*domain.User")).Return(nil)

	err := userService.RegisterUser(user)

	assert.NoError(t, err)
	mockRepo.AssertExpectations(t)
}
