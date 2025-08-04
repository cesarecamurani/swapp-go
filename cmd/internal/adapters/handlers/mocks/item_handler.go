package mocks

import (
	"github.com/google/uuid"
	"github.com/stretchr/testify/mock"
	"swapp-go/cmd/internal/domain"
)

type MockItemService struct {
	mock.Mock
}

func (m *MockItemService) Create(item *domain.Item) error {
	return m.Called(item).Error(0)
}

func (m *MockItemService) Update(id uuid.UUID, fields map[string]interface{}) (*domain.Item, error) {
	args := m.Called(id, fields)
	return args.Get(0).(*domain.Item), args.Error(1)
}

func (m *MockItemService) Delete(id uuid.UUID) error {
	return m.Called(id).Error(0)
}

func (m *MockItemService) FindByID(id uuid.UUID) (*domain.Item, error) {
	args := m.Called(id)

	if args.Get(0) == nil {
		return nil, args.Error(1)
	}

	return args.Get(0).(*domain.Item), args.Error(1)
}
