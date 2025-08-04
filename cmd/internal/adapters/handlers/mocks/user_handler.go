package mocks

import (
	"github.com/google/uuid"
	"github.com/stretchr/testify/mock"
	"swapp-go/cmd/internal/domain"
)

type MockUserService struct {
	mock.Mock
}

func (m *MockUserService) RegisterUser(user *domain.User) error {
	return m.Called(user).Error(0)
}

func (m *MockUserService) Update(id uuid.UUID, fields map[string]interface{}) (*domain.User, error) {
	args := m.Called(id, fields)

	if user, ok := args.Get(0).(*domain.User); ok {
		return user, args.Error(1)
	}

	return nil, args.Error(1)
}

func (m *MockUserService) Delete(id uuid.UUID) error {
	return m.Called(id).Error(0)
}

func (m *MockUserService) FindByID(id uuid.UUID) (*domain.User, error) {
	args := m.Called(id)

	if user, ok := args.Get(0).(*domain.User); ok {
		return user, args.Error(1)
	}

	return nil, args.Error(1)
}

func (m *MockUserService) FindByUsername(username string) (*domain.User, error) {
	args := m.Called(username)

	if user, ok := args.Get(0).(*domain.User); ok {
		return user, args.Error(1)
	}

	return nil, args.Error(1)
}

func (m *MockUserService) FindByEmail(email string) (*domain.User, error) {
	args := m.Called(email)

	if user, ok := args.Get(0).(*domain.User); ok {
		return user, args.Error(1)
	}

	return nil, args.Error(1)
}

func (m *MockUserService) Authenticate(username, password string) (string, *domain.User, error) {
	args := m.Called(username, password)
	token, _ := args.Get(0).(string)
	user, _ := args.Get(1).(*domain.User)

	return token, user, args.Error(2)
}
