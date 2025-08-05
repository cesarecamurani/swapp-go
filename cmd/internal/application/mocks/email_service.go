package mocks

import (
	"github.com/stretchr/testify/mock"
	"swapp-go/cmd/internal/domain"
)

type MockEmailService struct {
	mock.Mock
}

func (m *MockEmailService) SendEmail(message *domain.EmailMessage) error {
	return m.Called(message).Error(0)
}
