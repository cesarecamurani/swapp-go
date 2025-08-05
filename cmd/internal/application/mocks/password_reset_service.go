package mocks

import (
	"github.com/stretchr/testify/mock"
	"swapp-go/cmd/internal/domain"
)

type MockPasswordResetRepository struct {
	mock.Mock
}

func (m *MockPasswordResetRepository) Save(token *domain.PasswordReset) error {
	args := m.Called(token)
	return args.Error(0)
}

func (m *MockPasswordResetRepository) GetByToken(token string) (*domain.PasswordReset, error) {
	args := m.Called(token)
	result := args.Get(0)
	if result == nil {
		return nil, args.Error(1)
	}
	return result.(*domain.PasswordReset), args.Error(1)
}

func (m *MockPasswordResetRepository) Delete(token string) error {
	args := m.Called(token)
	return args.Error(0)
}
