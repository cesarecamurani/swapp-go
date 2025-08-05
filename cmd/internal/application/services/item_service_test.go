package services_test

import (
	"errors"
	"swapp-go/cmd/internal/application/mocks"
	"testing"

	"github.com/google/uuid"
	"github.com/stretchr/testify/assert"
	"swapp-go/cmd/internal/application/services"
	"swapp-go/cmd/internal/domain"
)

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

func TestItemService(t *testing.T) {
	t.Run("CreateItem_Success", func(t *testing.T) {
		mockRepo := new(mocks.ItemRepository)
		service := services.NewItemService(mockRepo)

		item := Item(uuid.New())
		mockRepo.On("Create", item).Return(nil)

		err := service.Create(item)
		assert.NoError(t, err)
		mockRepo.AssertExpectations(t)
	})

	t.Run("UpdateItem_Success", func(t *testing.T) {
		mockRepo := new(mocks.ItemRepository)
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
	})

	t.Run("UpdateItem_NotFound", func(t *testing.T) {
		mockRepo := new(mocks.ItemRepository)
		service := services.NewItemService(mockRepo)

		itemID := uuid.New()
		fields := map[string]interface{}{"name": "Doesn't matter"}

		mockRepo.On("FindByID", itemID).Return((*domain.Item)(nil), errors.New("not found"))

		item, err := service.Update(itemID, fields)
		assert.Nil(t, item)
		assert.Error(t, err)

		mockRepo.AssertExpectations(t)
	})

	t.Run("DeleteItem_Success", func(t *testing.T) {
		mockRepo := new(mocks.ItemRepository)
		service := services.NewItemService(mockRepo)

		itemID := uuid.New()
		mockRepo.On("Delete", itemID).Return(nil)

		err := service.Delete(itemID)
		assert.NoError(t, err)
		mockRepo.AssertExpectations(t)
	})

	t.Run("GetItemByID_Success", func(t *testing.T) {
		mockRepo := new(mocks.ItemRepository)
		service := services.NewItemService(mockRepo)

		itemID := uuid.New()
		item := Item(itemID)

		mockRepo.On("FindByID", itemID).Return(item, nil)

		result, err := service.FindByID(itemID)
		assert.NoError(t, err)
		assert.Equal(t, item, result)
		mockRepo.AssertExpectations(t)
	})
}
