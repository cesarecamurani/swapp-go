package services_test

import (
	"errors"
	"github.com/google/uuid"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
	"swapp-go/cmd/internal/application/services"
	"swapp-go/cmd/internal/domain"
	"testing"
)

// Mocks
type MockItemRepository struct {
	mock.Mock
}

func (m *MockItemRepository) Create(item *domain.Item) error {
	return m.Called(item).Error(0)
}

func (m *MockItemRepository) Update(id uuid.UUID, fields map[string]interface{}) (*domain.Item, error) {
	args := m.Called(id, fields)
	return args.Get(0).(*domain.Item), args.Error(1)
}

func (m *MockItemRepository) Delete(id uuid.UUID) error {
	return m.Called(id).Error(0)
}

func (m *MockItemRepository) FindByID(id uuid.UUID) (*domain.Item, error) {
	args := m.Called(id)
	return args.Get(0).(*domain.Item), args.Error(1)
}

// Test helpers
func Item(id uuid.UUID) *domain.Item {
	return &domain.Item{
		ID:          id,
		Name:        "Test",
		Description: "Sample",
		PictureURL:  "/uploads/test.jpg",
		UserID:      uuid.New(),
	}
}

// Tests
func TestCreateItem_Success(t *testing.T) {
	mockRepo := new(MockItemRepository)
	service := services.NewItemService(mockRepo)

	item := Item(uuid.New())
	mockRepo.On("Create", item).Return(nil)

	err := service.Create(item)
	assert.NoError(t, err)
	mockRepo.AssertExpectations(t)
}

func TestUpdateItem_Success(t *testing.T) {
	mockRepo := new(MockItemRepository)
	service := services.NewItemService(mockRepo)

	itemID := uuid.New()
	fields := map[string]interface{}{"name": "Updated"}
	item := Item(itemID)

	mockRepo.On("FindByID", itemID).Return(item, nil)
	mockRepo.On("Update", itemID, fields).Return(item, nil)

	updated, err := service.Update(itemID, fields)
	assert.NoError(t, err)
	assert.Equal(t, item, updated)

	mockRepo.AssertExpectations(t)
}

func TestUpdateItem_NotFound(t *testing.T) {
	mockRepo := new(MockItemRepository)
	service := services.NewItemService(mockRepo)

	itemID := uuid.New()
	fields := map[string]interface{}{"name": "Doesn't matter"}

	mockRepo.
		On("FindByID", itemID).
		Return((*domain.Item)(nil), errors.New("not found"))

	item, err := service.Update(itemID, fields)
	assert.Nil(t, item)
	assert.Error(t, err)

	mockRepo.AssertExpectations(t)
}

func TestDeleteItem_Success(t *testing.T) {
	mockRepo := new(MockItemRepository)
	service := services.NewItemService(mockRepo)

	itemID := uuid.New()
	mockRepo.On("Delete", itemID).Return(nil)

	err := service.Delete(itemID)
	assert.NoError(t, err)
	mockRepo.AssertExpectations(t)
}

func TestGetItemByID_Success(t *testing.T) {
	mockRepo := new(MockItemRepository)
	service := services.NewItemService(mockRepo)

	itemID := uuid.New()
	item := Item(itemID)

	mockRepo.On("FindByID", itemID).Return(item, nil)

	result, err := service.FindByID(itemID)
	assert.NoError(t, err)
	assert.Equal(t, item, result)
	mockRepo.AssertExpectations(t)
}
