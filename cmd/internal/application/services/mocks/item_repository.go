package mocks

import (
	"github.com/google/uuid"
	"github.com/stretchr/testify/mock"
	"swapp-go/cmd/internal/domain"
)

type ItemRepository struct {
	mock.Mock
}

func (m *ItemRepository) Create(item *domain.Item) error {
	args := m.Called(item)
	return args.Error(0)
}

func (m *ItemRepository) FindByID(id uuid.UUID) (*domain.Item, error) {
	args := m.Called(id)
	if item, ok := args.Get(0).(*domain.Item); ok {
		return item, args.Error(1)
	}
	return nil, args.Error(1)
}

func (m *ItemRepository) Update(id uuid.UUID, fields map[string]interface{}) (*domain.Item, error) {
	args := m.Called(id, fields)
	if item, ok := args.Get(0).(*domain.Item); ok {
		return item, args.Error(1)
	}
	return nil, args.Error(1)
}

func (m *ItemRepository) Delete(id uuid.UUID) error {
	args := m.Called(id)
	return args.Error(0)
}

func (m *ItemRepository) TryMarkItemAsOffered(itemID uuid.UUID) (bool, error) {
	args := m.Called(itemID)
	return args.Bool(0), args.Error(1)
}
