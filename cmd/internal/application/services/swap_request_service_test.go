package services_test

import (
	"errors"
	"github.com/google/uuid"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
	"swapp-go/cmd/internal/application/services"
	"swapp-go/cmd/internal/application/services/mocks"
	"swapp-go/cmd/internal/domain"
	"testing"
)

func setupSwapRequestServiceTest() (*services.SwapRequestService, *mocks.SwapRequestRepository, *mocks.ItemRepository) {
	mockSwapRequestRepo := new(mocks.SwapRequestRepository)
	mockItemRepo := new(mocks.ItemRepository)
	service := services.NewSwapRequestService(mockSwapRequestRepo, mockItemRepo)
	return service, mockSwapRequestRepo, mockItemRepo
}

func TestSwapRequestService_Create(t *testing.T) {
	testItemID := uuid.New()
	testRequest := &domain.SwapRequest{
		OfferedItemID: testItemID,
	}

	t.Run("success", func(t *testing.T) {
		service, mockSwapRequestRepo, mockItemRepo := setupSwapRequestServiceTest()
		mockItemRepo.On("FindByID", testItemID).Return(&domain.Item{ID: testItemID}, nil).Once()
		mockItemRepo.On("TryMarkItemAsOffered", testItemID).Return(true, nil).Once()
		mockItemRepo.On("Update", testItemID, mock.Anything).Return(&domain.Item{ID: testItemID}, nil).Once()
		mockSwapRequestRepo.On("Create", testRequest).Return(nil).Once()

		err := service.Create(testRequest)
		assert.NoError(t, err)

		mockItemRepo.AssertExpectations(t)
		mockSwapRequestRepo.AssertExpectations(t)
	})

	t.Run("offered item not found", func(t *testing.T) {
		service, _, mockItemRepo := setupSwapRequestServiceTest()
		mockItemRepo.On("FindByID", testItemID).Return(nil, errors.New("not found")).Once()

		err := service.Create(testRequest)
		assert.EqualError(t, err, "offered item not found")
		mockItemRepo.AssertExpectations(t)
	})

	t.Run("item already offered", func(t *testing.T) {
		service, _, mockItemRepo := setupSwapRequestServiceTest()
		mockItemRepo.On("FindByID", testItemID).Return(&domain.Item{ID: testItemID}, nil).Once()
		mockItemRepo.On("TryMarkItemAsOffered", testItemID).Return(false, nil).Once()

		err := service.Create(testRequest)
		assert.ErrorIs(t, err, services.ItemAlreadyOfferedErr)
		mockItemRepo.AssertExpectations(t)
	})

	t.Run("TryMarkItemAsOffered error", func(t *testing.T) {
		service, _, mockItemRepo := setupSwapRequestServiceTest()
		mockItemRepo.On("FindByID", testItemID).Return(&domain.Item{ID: testItemID}, nil).Once()
		mockItemRepo.On("TryMarkItemAsOffered", testItemID).Return(false, errors.New("db error")).Once()

		err := service.Create(testRequest)
		assert.EqualError(t, err, "db error")
		mockItemRepo.AssertExpectations(t)
	})

	t.Run("repo.Create error rolls back offered flag", func(t *testing.T) {
		service, mockSwapRequestRepo, mockItemRepo := setupSwapRequestServiceTest()
		mockItemRepo.On("FindByID", testItemID).Return(&domain.Item{ID: testItemID}, nil).Once()
		mockItemRepo.On("TryMarkItemAsOffered", testItemID).Return(true, nil).Once()
		mockItemRepo.On("Update", testItemID, mock.MatchedBy(func(m map[string]interface{}) bool {
			return m["offered"] == true
		})).Return(&domain.Item{ID: testItemID}, nil).Once()
		mockSwapRequestRepo.On("Create", testRequest).Return(errors.New("create error")).Once()
		mockItemRepo.On("Update", testItemID, mock.MatchedBy(func(m map[string]interface{}) bool {
			return m["offered"] == false
		})).Return(&domain.Item{ID: testItemID}, nil).Once()

		err := service.Create(testRequest)
		assert.EqualError(t, err, "create error")
		mockItemRepo.AssertExpectations(t)
		mockSwapRequestRepo.AssertExpectations(t)
	})
}

func TestSwapRequestService_FindByID(t *testing.T) {
	existingID := uuid.New()
	expectedSwapRequest := &domain.SwapRequest{ID: existingID}

	t.Run("success", func(t *testing.T) {
		service, mockSwapRequestRepo, _ := setupSwapRequestServiceTest()
		mockSwapRequestRepo.On("FindByID", existingID).Return(expectedSwapRequest, nil).Once()

		result, err := service.FindByID(existingID)
		assert.NoError(t, err)
		assert.Equal(t, expectedSwapRequest, result)

		mockSwapRequestRepo.AssertExpectations(t)
	})

	t.Run("not found error", func(t *testing.T) {
		service, mockSwapRequestRepo, _ := setupSwapRequestServiceTest()
		mockSwapRequestRepo.On("FindByID", existingID).Return(nil, errors.New("not found")).Once()

		result, err := service.FindByID(existingID)
		assert.Error(t, err)
		assert.Nil(t, result)

		mockSwapRequestRepo.AssertExpectations(t)
	})
}

func TestSwapRequestService_FindByReferenceNumber(t *testing.T) {
	reference := "REF123"
	expected := &domain.SwapRequest{ReferenceNumber: reference}

	t.Run("success", func(t *testing.T) {
		service, mockSwapRequestRepo, _ := setupSwapRequestServiceTest()
		mockSwapRequestRepo.On("FindByReferenceNumber", reference).Return(expected, nil).Once()

		result, err := service.FindByReferenceNumber(reference)
		assert.NoError(t, err)
		assert.Equal(t, expected, result)

		mockSwapRequestRepo.AssertExpectations(t)
	})

	t.Run("not found error", func(t *testing.T) {
		service, mockSwapRequestRepo, _ := setupSwapRequestServiceTest()
		mockSwapRequestRepo.On("FindByReferenceNumber", reference).Return(nil, errors.New("not found")).Once()

		result, err := service.FindByReferenceNumber(reference)
		assert.Error(t, err)
		assert.Nil(t, result)

		mockSwapRequestRepo.AssertExpectations(t)
	})
}

func TestSwapRequestService_ListByUser(t *testing.T) {
	userID := uuid.New()
	expectedList := []domain.SwapRequest{{ID: uuid.New()}, {ID: uuid.New()}}

	t.Run("success", func(t *testing.T) {
		service, mockSwapRequestRepo, _ := setupSwapRequestServiceTest()
		mockSwapRequestRepo.On("ListByUser", userID).Return(expectedList, nil).Once()

		result, err := service.ListByUser(userID)
		assert.NoError(t, err)
		assert.Equal(t, expectedList, result)

		mockSwapRequestRepo.AssertExpectations(t)
	})

	t.Run("error from repo", func(t *testing.T) {
		service, mockSwapRequestRepo, _ := setupSwapRequestServiceTest()
		mockSwapRequestRepo.On("ListByUser", userID).Return(nil, errors.New("db error")).Once()

		result, err := service.ListByUser(userID)
		assert.Error(t, err)
		assert.Nil(t, result)

		mockSwapRequestRepo.AssertExpectations(t)
	})
}

func TestSwapRequestService_ListByStatus(t *testing.T) {
	status := domain.StatusPending
	expectedList := []domain.SwapRequest{
		{ID: uuid.New(), Status: status},
		{ID: uuid.New(), Status: status},
	}

	t.Run("success", func(t *testing.T) {
		service, mockSwapRequestRepo, _ := setupSwapRequestServiceTest()
		mockSwapRequestRepo.On("ListByStatus", status).Return(expectedList, nil).Once()

		result, err := service.ListByStatus(status)
		assert.NoError(t, err)
		assert.Equal(t, expectedList, result)

		mockSwapRequestRepo.AssertExpectations(t)
	})

	t.Run("error from repo", func(t *testing.T) {
		service, mockSwapRequestRepo, _ := setupSwapRequestServiceTest()
		mockSwapRequestRepo.On("ListByStatus", status).Return(nil, errors.New("db error")).Once()

		result, err := service.ListByStatus(status)
		assert.Error(t, err)
		assert.Nil(t, result)

		mockSwapRequestRepo.AssertExpectations(t)
	})
}

func TestSwapRequestService_UpdateStatus(t *testing.T) {
	swapRequestID := uuid.New()
	offeredItemID := uuid.New()
	swapRequest := &domain.SwapRequest{ID: swapRequestID, OfferedItemID: offeredItemID}

	t.Run("success - not rejected/cancelled", func(t *testing.T) {
		service, mockSwapRequestRepo, mockItemRepo := setupSwapRequestServiceTest()
		mockSwapRequestRepo.On("FindByID", swapRequestID).Return(swapRequest, nil).Once()
		mockSwapRequestRepo.On("UpdateStatus", swapRequestID, domain.StatusAccepted).Return(nil).Once()

		err := service.UpdateStatus(swapRequestID, domain.StatusAccepted)
		assert.NoError(t, err)

		mockItemRepo.AssertNotCalled(t, "Update", mock.Anything, mock.Anything)
		mockSwapRequestRepo.AssertExpectations(t)
	})

	t.Run("success - rejected", func(t *testing.T) {
		service, mockSwapRequestRepo, mockItemRepo := setupSwapRequestServiceTest()
		mockSwapRequestRepo.On("FindByID", swapRequestID).Return(swapRequest, nil).Once()
		mockSwapRequestRepo.On("UpdateStatus", swapRequestID, domain.StatusRejected).Return(nil).Once()
		mockItemRepo.On("Update", offeredItemID, mock.MatchedBy(func(fields map[string]interface{}) bool {
			return fields["offered"] == false
		})).Return(&domain.Item{}, nil).Once()

		err := service.UpdateStatus(swapRequestID, domain.StatusRejected)
		assert.NoError(t, err)

		mockItemRepo.AssertExpectations(t)
		mockSwapRequestRepo.AssertExpectations(t)
	})

	t.Run("success - cancelled", func(t *testing.T) {
		service, mockSwapRequestRepo, mockItemRepo := setupSwapRequestServiceTest()
		mockSwapRequestRepo.On("FindByID", swapRequestID).Return(swapRequest, nil).Once()
		mockSwapRequestRepo.On("UpdateStatus", swapRequestID, domain.StatusCancelled).Return(nil).Once()
		mockItemRepo.On("Update", offeredItemID, mock.MatchedBy(func(fields map[string]interface{}) bool {
			return fields["offered"] == false
		})).Return(&domain.Item{}, nil).Once()

		err := service.UpdateStatus(swapRequestID, domain.StatusCancelled)
		assert.NoError(t, err)

		mockItemRepo.AssertExpectations(t)
		mockSwapRequestRepo.AssertExpectations(t)
	})

	t.Run("error - not found", func(t *testing.T) {
		service, mockSwapRequestRepo, mockItemRepo := setupSwapRequestServiceTest()
		mockSwapRequestRepo.On("FindByID", swapRequestID).Return(nil, errors.New("not found")).Once()

		err := service.UpdateStatus(swapRequestID, domain.StatusAccepted)
		assert.Error(t, err)

		mockItemRepo.AssertNotCalled(t, "Update", mock.Anything, mock.Anything)
		mockSwapRequestRepo.AssertExpectations(t)
	})

	t.Run("error - update failed", func(t *testing.T) {
		service, mockSwapRequestRepo, mockItemRepo := setupSwapRequestServiceTest()
		mockSwapRequestRepo.On("FindByID", swapRequestID).Return(swapRequest, nil).Once()
		mockSwapRequestRepo.On("UpdateStatus", swapRequestID, domain.StatusAccepted).Return(errors.New("update error")).Once()

		err := service.UpdateStatus(swapRequestID, domain.StatusAccepted)
		assert.Error(t, err)

		mockItemRepo.AssertNotCalled(t, "Update", mock.Anything, mock.Anything)
		mockSwapRequestRepo.AssertExpectations(t)
	})

	t.Run("error - reset offered status failed", func(t *testing.T) {
		service, mockSwapRequestRepo, mockItemRepo := setupSwapRequestServiceTest()
		mockSwapRequestRepo.On("FindByID", swapRequestID).Return(swapRequest, nil).Once()
		mockSwapRequestRepo.On("UpdateStatus", swapRequestID, domain.StatusRejected).Return(nil).Once()
		mockItemRepo.On("Update", offeredItemID, mock.Anything).Return(nil, errors.New("update error")).Once()

		err := service.UpdateStatus(swapRequestID, domain.StatusRejected)
		assert.Error(t, err)

		mockItemRepo.AssertExpectations(t)
		mockSwapRequestRepo.AssertExpectations(t)
	})
}

func TestSwapRequestService_Delete(t *testing.T) {
	swapRequestID := uuid.New()
	offeredItemID := uuid.New()
	swapRequest := &domain.SwapRequest{ID: swapRequestID, OfferedItemID: offeredItemID}

	t.Run("success", func(t *testing.T) {
		service, mockSwapRequestRepo, mockItemRepo := setupSwapRequestServiceTest()
		mockSwapRequestRepo.On("FindByID", swapRequestID).Return(swapRequest, nil).Once()
		mockSwapRequestRepo.On("Delete", swapRequestID).Return(nil).Once()
		mockItemRepo.On("Update", offeredItemID, mock.Anything).Return(&domain.Item{}, nil).Once()

		err := service.Delete(swapRequestID)
		assert.NoError(t, err)

		mockItemRepo.AssertExpectations(t)
		mockSwapRequestRepo.AssertExpectations(t)
	})

	t.Run("not found", func(t *testing.T) {
		service, mockSwapRequestRepo, mockItemRepo := setupSwapRequestServiceTest()
		mockSwapRequestRepo.On("FindByID", swapRequestID).Return(nil, errors.New("not found")).Once()

		err := service.Delete(swapRequestID)
		assert.Error(t, err)

		mockItemRepo.AssertNotCalled(t, "Update", mock.Anything, mock.Anything)
		mockSwapRequestRepo.AssertExpectations(t)
	})

	t.Run("delete failed", func(t *testing.T) {
		service, mockSwapRequestRepo, mockItemRepo := setupSwapRequestServiceTest()
		mockSwapRequestRepo.On("FindByID", swapRequestID).Return(swapRequest, nil).Once()
		mockSwapRequestRepo.On("Delete", swapRequestID).Return(errors.New("delete error")).Once()

		err := service.Delete(swapRequestID)
		assert.Error(t, err)

		mockItemRepo.AssertNotCalled(t, "Update", mock.Anything, mock.Anything)
		mockSwapRequestRepo.AssertExpectations(t)
	})

	t.Run("reset offered status failed", func(t *testing.T) {
		service, mockSwapRequestRepo, mockItemRepo := setupSwapRequestServiceTest()
		mockSwapRequestRepo.On("FindByID", swapRequestID).Return(swapRequest, nil).Once()
		mockSwapRequestRepo.On("Delete", swapRequestID).Return(nil).Once()
		mockItemRepo.On("Update", offeredItemID, mock.Anything).Return(nil, errors.New("update error")).Once()

		err := service.Delete(swapRequestID)
		assert.Error(t, err)

		mockItemRepo.AssertExpectations(t)
		mockSwapRequestRepo.AssertExpectations(t)
	})
}
