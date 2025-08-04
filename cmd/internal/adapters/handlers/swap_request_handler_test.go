package handlers_test

import (
	"bytes"
	"encoding/json"
	"errors"
	"net/http"
	"net/http/httptest"
	"swapp-go/cmd/internal/application/services"
	"testing"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"

	"swapp-go/cmd/internal/adapters/handlers"
	"swapp-go/cmd/internal/adapters/handlers/mocks"
	"swapp-go/cmd/internal/domain"
)

var (
	testUserID          = uuid.New()
	testSwapRequestID   = uuid.New()
	testReference       = "ref-xyz-123"
	testOfferedItemID   = uuid.New()
	testRequestedItemID = uuid.New()
	testRecipientID     = uuid.New()
)

func setupRouterAndHandler(mockService *mocks.SwapRequestService) *gin.Engine {
	gin.SetMode(gin.TestMode)
	router := gin.New() // don't use Default() to avoid logging during tests
	handler := handlers.NewSwapRequestHandler(mockService)

	// ✅ Middleware to inject testUserID into context
	router.Use(func(c *gin.Context) {
		c.Set("userID", testUserID.String())
		c.Next()
	})

	router.POST("/swap-requests/create", handler.Create)
	router.GET("/swap-requests/:id", handler.FindByID)
	router.GET("/swap-requests/reference/:reference", handler.FindByReferenceNumber)
	router.GET("/swap-requests/list-by-user/:id", handler.ListByUser)
	router.GET("/swap-requests/list-by-status/:status", handler.ListByStatus)
	router.PATCH("/swap-requests/update-status/:id", handler.UpdateStatus)
	router.DELETE("/swap-requests/delete/:id", handler.Delete)

	return router
}

func newTestRouter() (*gin.Engine, *mocks.SwapRequestService) {
	mockService := new(mocks.SwapRequestService)
	router := setupRouterAndHandler(mockService)

	return router, mockService
}

func TestSwapRequestHandler_Create(t *testing.T) {
	t.Run("success", func(t *testing.T) {
		router, mockService := newTestRouter()

		reqBody := handlers.SwapRequestRequest{
			OfferedItemID:   testOfferedItemID,
			RequestedItemID: testRequestedItemID,
			RecipientID:     testRecipientID,
		}
		jsonBody, _ := json.Marshal(reqBody)

		mockService.On("Create", mock.AnythingOfType("*domain.SwapRequest")).Return(nil).Run(func(args mock.Arguments) {
			arg := args.Get(0).(*domain.SwapRequest)
			arg.ID = testSwapRequestID
		})

		req := httptest.NewRequest(http.MethodPost, "/swap-requests/create", bytes.NewBuffer(jsonBody))
		req.Header.Set("Content-Type", "application/json")
		resp := httptest.NewRecorder()

		router.ServeHTTP(resp, req)
		assert.Equal(t, http.StatusCreated, resp.Code)
		mockService.AssertExpectations(t)
	})

	t.Run("invalid JSON", func(t *testing.T) {
		router, _ := newTestRouter()

		req := httptest.NewRequest(http.MethodPost, "/swap-requests/create", bytes.NewBufferString(`{"offered_item_id":}`))
		req.Header.Set("Content-Type", "application/json")
		resp := httptest.NewRecorder()

		router.ServeHTTP(resp, req)
		assert.Equal(t, http.StatusBadRequest, resp.Code)
	})

	t.Run("item already offered", func(t *testing.T) {
		router, mockService := newTestRouter()

		reqBody := handlers.SwapRequestRequest{
			OfferedItemID:   testOfferedItemID,
			RequestedItemID: testRequestedItemID,
			RecipientID:     testRecipientID,
		}
		jsonBody, _ := json.Marshal(reqBody)

		mockService.
			On("Create", mock.AnythingOfType("*domain.SwapRequest")).
			Return(services.ItemAlreadyOfferedErr)

		req := httptest.NewRequest(http.MethodPost, "/swap-requests/create", bytes.NewBuffer(jsonBody))
		req.Header.Set("Content-Type", "application/json")
		resp := httptest.NewRecorder()

		router.ServeHTTP(resp, req)
		assert.Equal(t, http.StatusConflict, resp.Code)
	})

	t.Run("internal error", func(t *testing.T) {
		router, mockService := newTestRouter()

		reqBody := handlers.SwapRequestRequest{
			OfferedItemID:   testOfferedItemID,
			RequestedItemID: testRequestedItemID,
			RecipientID:     testRecipientID,
		}
		jsonBody, _ := json.Marshal(reqBody)

		mockService.On("Create", mock.AnythingOfType("*domain.SwapRequest")).Return(errors.New("fail"))

		req := httptest.NewRequest(http.MethodPost, "/swap-requests/create", bytes.NewBuffer(jsonBody))
		req.Header.Set("Content-Type", "application/json")
		resp := httptest.NewRecorder()

		router.ServeHTTP(resp, req)
		assert.Equal(t, http.StatusInternalServerError, resp.Code)
	})
}

func TestSwapRequestHandler_FindByID(t *testing.T) {
	swapRequest := &domain.SwapRequest{
		ID:              testSwapRequestID,
		SenderID:        testUserID,
		RecipientID:     testRecipientID,
		OfferedItemID:   testOfferedItemID,
		RequestedItemID: testRequestedItemID,
		Status:          domain.StatusPending,
		ReferenceNumber: testReference,
	}

	t.Run("success", func(t *testing.T) {
		router, mockService := newTestRouter()

		mockService.On("FindByID", testSwapRequestID).Return(swapRequest, nil)

		req := httptest.NewRequest(http.MethodGet, "/swap-requests/"+testSwapRequestID.String(), nil)
		resp := httptest.NewRecorder()

		router.ServeHTTP(resp, req)
		assert.Equal(t, http.StatusOK, resp.Code)
		mockService.AssertCalled(t, "FindByID", testSwapRequestID)
	})

	t.Run("not found", func(t *testing.T) {
		router, mockService := newTestRouter()

		mockService.On("FindByID", testSwapRequestID).Return(nil, errors.New("not found"))

		req := httptest.NewRequest(http.MethodGet, "/swap-requests/"+testSwapRequestID.String(), nil)
		resp := httptest.NewRecorder()

		router.ServeHTTP(resp, req)
		assert.Equal(t, http.StatusNotFound, resp.Code)
	})

	t.Run("unauthorized", func(t *testing.T) {
		router, mockService := newTestRouter()

		unauthorizedRequest := &domain.SwapRequest{
			ID:          testSwapRequestID,
			SenderID:    uuid.New(),
			RecipientID: uuid.New(),
			Status:      domain.StatusPending,
		}
		mockService.On("FindByID", testSwapRequestID).Return(unauthorizedRequest, nil)

		req := httptest.NewRequest(http.MethodGet, "/swap-requests/"+testSwapRequestID.String(), nil)
		resp := httptest.NewRecorder()

		router.ServeHTTP(resp, req)
		assert.Equal(t, http.StatusUnauthorized, resp.Code)
	})
}

func TestSwapRequestHandler_FindByReferenceNumber(t *testing.T) {
	swapRequest := &domain.SwapRequest{
		ID:              testSwapRequestID,
		SenderID:        testUserID,
		RecipientID:     testRecipientID,
		Status:          domain.StatusPending,
		ReferenceNumber: testReference,
		OfferedItemID:   testOfferedItemID,
		RequestedItemID: testRequestedItemID,
	}

	t.Run("success", func(t *testing.T) {
		router, mockService := newTestRouter()

		mockService.On("FindByReferenceNumber", testReference).Return(swapRequest, nil)

		req := httptest.NewRequest(http.MethodGet, "/swap-requests/reference/"+testReference, nil)
		resp := httptest.NewRecorder()

		router.ServeHTTP(resp, req)
		assert.Equal(t, http.StatusOK, resp.Code)
	})

	t.Run("not found", func(t *testing.T) {
		router, mockService := newTestRouter()

		mockService.On("FindByReferenceNumber", testReference).Return(nil, errors.New("not found"))

		req := httptest.NewRequest(http.MethodGet, "/swap-requests/reference/"+testReference, nil)
		resp := httptest.NewRecorder()

		router.ServeHTTP(resp, req)
		assert.Equal(t, http.StatusNotFound, resp.Code)
	})

	t.Run("unauthorized", func(t *testing.T) {
		router, mockService := newTestRouter()

		unauthorizedSwap := &domain.SwapRequest{
			SenderID:    uuid.New(),
			RecipientID: uuid.New(),
		}
		mockService.On("FindByReferenceNumber", "ref-unauthorized").Return(unauthorizedSwap, nil)

		req := httptest.NewRequest(http.MethodGet, "/swap-requests/reference/ref-unauthorized", nil)
		resp := httptest.NewRecorder()

		router.ServeHTTP(resp, req)
		assert.Equal(t, http.StatusUnauthorized, resp.Code)
	})
}

func TestSwapRequestHandler_ListByUser(t *testing.T) {
	swapRequests := []domain.SwapRequest{
		{
			ID:          testSwapRequestID,
			SenderID:    testUserID,
			RecipientID: testRecipientID,
		},
	}

	t.Run("success", func(t *testing.T) {
		router, mockService := newTestRouter()

		mockService.On("ListByUser", testUserID).Return(swapRequests, nil)

		req := httptest.NewRequest(http.MethodGet, "/swap-requests/list-by-user/"+testUserID.String(), nil)
		resp := httptest.NewRecorder()

		router.ServeHTTP(resp, req)
		assert.Equal(t, http.StatusOK, resp.Code)
	})

	t.Run("service error", func(t *testing.T) {
		router, mockService := newTestRouter()

		mockService.On("ListByUser", testUserID).Return(nil, errors.New("fail"))

		req := httptest.NewRequest(http.MethodGet, "/swap-requests/list-by-user/"+testUserID.String(), nil)
		resp := httptest.NewRecorder()

		router.ServeHTTP(resp, req)
		assert.Equal(t, http.StatusInternalServerError, resp.Code)
	})
}

func TestSwapRequestHandler_ListByStatus(t *testing.T) {
	status := domain.StatusPending
	swapRequests := []domain.SwapRequest{
		{ID: testSwapRequestID, SenderID: testUserID, RecipientID: uuid.New(), Status: status},
	}

	t.Run("success with valid status", func(t *testing.T) {
		router, mockService := newTestRouter()

		mockService.On("ListByStatus", status).Return(swapRequests, nil)

		req := httptest.NewRequest(http.MethodGet, "/swap-requests/list-by-status/"+string(status), nil)
		resp := httptest.NewRecorder()

		router.ServeHTTP(resp, req)
		assert.Equal(t, http.StatusOK, resp.Code)
	})

	t.Run("invalid status param", func(t *testing.T) {
		router, _ := newTestRouter()

		req := httptest.NewRequest(http.MethodGet, "/swap-requests/list-by-status/invalid-status", nil)
		resp := httptest.NewRecorder()

		router.ServeHTTP(resp, req)
		assert.Equal(t, http.StatusBadRequest, resp.Code)
	})

	t.Run("service error", func(t *testing.T) {
		router, mockService := newTestRouter()

		mockService.On("ListByStatus", domain.StatusRejected).Return(nil, errors.New("fail"))

		req := httptest.NewRequest(http.MethodGet, "/swap-requests/list-by-status/rejected", nil)
		resp := httptest.NewRecorder()

		router.ServeHTTP(resp, req)
		assert.Equal(t, http.StatusInternalServerError, resp.Code)
	})
}

func TestSwapRequestHandler_Delete(t *testing.T) {
	swapRequest := &domain.SwapRequest{
		ID:          testSwapRequestID,
		SenderID:    testUserID,
		RecipientID: testRecipientID,
	}

	t.Run("success", func(t *testing.T) {
		router, mockService := newTestRouter()

		mockService.On("FindByID", testSwapRequestID).Return(swapRequest, nil)
		mockService.On("Delete", testSwapRequestID).Return(nil)

		req := httptest.NewRequest(http.MethodDelete, "/swap-requests/delete/"+testSwapRequestID.String(), nil)
		resp := httptest.NewRecorder()

		router.ServeHTTP(resp, req)
		assert.Equal(t, http.StatusOK, resp.Code)
	})

	t.Run("unauthorized", func(t *testing.T) {
		router, mockService := newTestRouter()

		unauthorizedSwapRequest := &domain.SwapRequest{
			ID:          testSwapRequestID,
			SenderID:    uuid.New(),
			RecipientID: uuid.New(),
		}
		mockService.On("FindByID", testSwapRequestID).Return(unauthorizedSwapRequest, nil)

		req := httptest.NewRequest(http.MethodDelete, "/swap-requests/delete/"+testSwapRequestID.String(), nil)
		resp := httptest.NewRecorder()

		router.ServeHTTP(resp, req)
		assert.Equal(t, http.StatusUnauthorized, resp.Code)
	})

	t.Run("not found", func(t *testing.T) {
		router, mockService := newTestRouter()

		mockService.On("FindByID", testSwapRequestID).Return(nil, errors.New("not found"))

		req := httptest.NewRequest(http.MethodDelete, "/swap-requests/delete/"+testSwapRequestID.String(), nil)
		resp := httptest.NewRecorder()

		router.ServeHTTP(resp, req)
		assert.Equal(t, http.StatusNotFound, resp.Code)
	})

	t.Run("delete error", func(t *testing.T) {
		router, mockService := newTestRouter()

		mockService.On("FindByID", testSwapRequestID).Return(swapRequest, nil)
		mockService.On("Delete", testSwapRequestID).Return(errors.New("fail"))

		req := httptest.NewRequest(http.MethodDelete, "/swap-requests/delete/"+testSwapRequestID.String(), nil)
		resp := httptest.NewRecorder()

		router.ServeHTTP(resp, req)
		assert.Equal(t, http.StatusInternalServerError, resp.Code)
	})
}

func TestSwapRequestHandler_UpdateStatus(t *testing.T) {
	swapID := uuid.New()

	t.Run("success_cancel_by_sender", func(t *testing.T) {
		router, mockService := newTestRouter()

		swap := &domain.SwapRequest{
			ID:              swapID,
			SenderID:        testUserID,
			RecipientID:     uuid.New(),
			Status:          domain.StatusPending,
			OfferedItemID:   uuid.New(),
			RequestedItemID: uuid.New(),
		}

		mockService.On("FindByID", swapID).Return(swap, nil)
		mockService.On("UpdateStatus", swapID, domain.StatusCancelled).Return(nil)

		body := map[string]string{"status": "cancelled"}
		jsonBody, _ := json.Marshal(body)

		req := httptest.NewRequest(http.MethodPatch, "/swap-requests/update-status/"+swapID.String(), bytes.NewBuffer(jsonBody))
		req.Header.Set("Content-Type", "application/json")
		resp := httptest.NewRecorder()

		router.ServeHTTP(resp, req)
		assert.Equal(t, http.StatusOK, resp.Code)
	})

	t.Run("success_accept_by_recipient", func(t *testing.T) {
		router, mockService := newTestRouter()

		swap := &domain.SwapRequest{
			ID:              swapID,
			SenderID:        uuid.New(),
			RecipientID:     testUserID,
			Status:          domain.StatusPending,
			OfferedItemID:   uuid.New(),
			RequestedItemID: uuid.New(),
		}

		mockService.On("FindByID", swapID).Return(swap, nil)
		mockService.On("UpdateStatus", swapID, domain.StatusAccepted).Return(nil)

		body := map[string]string{"status": "accepted"}
		jsonBody, _ := json.Marshal(body)

		req := httptest.NewRequest(http.MethodPatch, "/swap-requests/update-status/"+swapID.String(), bytes.NewBuffer(jsonBody))
		req.Header.Set("Content-Type", "application/json")
		resp := httptest.NewRecorder()

		router.ServeHTTP(resp, req)
		assert.Equal(t, http.StatusOK, resp.Code)
	})

	t.Run("success_reject_by_recipient", func(t *testing.T) {
		router, mockService := newTestRouter()

		swap := &domain.SwapRequest{
			ID:              swapID,
			SenderID:        uuid.New(),
			RecipientID:     testUserID,
			Status:          domain.StatusPending,
			OfferedItemID:   uuid.New(),
			RequestedItemID: uuid.New(),
		}

		mockService.On("FindByID", swapID).Return(swap, nil)
		mockService.On("UpdateStatus", swapID, domain.StatusRejected).Return(nil)

		body := map[string]string{"status": "rejected"}
		jsonBody, _ := json.Marshal(body)

		req := httptest.NewRequest(http.MethodPatch, "/swap-requests/update-status/"+swapID.String(), bytes.NewBuffer(jsonBody))
		req.Header.Set("Content-Type", "application/json")
		resp := httptest.NewRecorder()

		router.ServeHTTP(resp, req)
		assert.Equal(t, http.StatusOK, resp.Code)
	})

	t.Run("unauthorized_user", func(t *testing.T) {
		router, mockService := newTestRouter()

		swap := &domain.SwapRequest{
			ID:              swapID,
			SenderID:        uuid.New(),
			RecipientID:     uuid.New(),
			Status:          domain.StatusPending,
			OfferedItemID:   uuid.New(),
			RequestedItemID: uuid.New(),
		}

		mockService.On("FindByID", swapID).Return(swap, nil)

		body := map[string]string{"status": "accepted"}
		jsonBody, _ := json.Marshal(body)

		req := httptest.NewRequest(http.MethodPatch, "/swap-requests/update-status/"+swapID.String(), bytes.NewBuffer(jsonBody))
		req.Header.Set("Content-Type", "application/json")
		resp := httptest.NewRecorder()

		router.ServeHTTP(resp, req)
		assert.Equal(t, http.StatusUnauthorized, resp.Code)
	})

	t.Run("invalid_JSON_payload", func(t *testing.T) {
		router, _ := newTestRouter()

		req := httptest.NewRequest(http.MethodPatch, "/swap-requests/update-status/"+swapID.String(), bytes.NewBuffer([]byte(`invalid-json`)))
		req.Header.Set("Content-Type", "application/json")
		resp := httptest.NewRecorder()

		router.ServeHTTP(resp, req)
		assert.Equal(t, http.StatusBadRequest, resp.Code)
	})

	t.Run("invalid_status_transition", func(t *testing.T) {
		router, mockService := newTestRouter()

		swap := &domain.SwapRequest{
			ID:              swapID,
			SenderID:        testUserID,
			RecipientID:     uuid.New(),
			Status:          domain.StatusPending,
			OfferedItemID:   uuid.New(),
			RequestedItemID: uuid.New(),
		}

		mockService.On("FindByID", swapID).Return(swap, nil)

		body := map[string]string{"status": "not-a-valid-status"}
		jsonBody, _ := json.Marshal(body)

		req := httptest.NewRequest(http.MethodPatch, "/swap-requests/update-status/"+swapID.String(), bytes.NewBuffer(jsonBody))
		req.Header.Set("Content-Type", "application/json")
		resp := httptest.NewRecorder()

		router.ServeHTTP(resp, req)
		assert.Equal(t, http.StatusBadRequest, resp.Code)
	})

	t.Run("update_status_service_error", func(t *testing.T) {
		router, mockService := newTestRouter()

		swap := &domain.SwapRequest{
			ID:              swapID,
			SenderID:        uuid.New(),
			RecipientID:     testUserID,
			Status:          domain.StatusPending,
			OfferedItemID:   uuid.New(),
			RequestedItemID: uuid.New(),
		}

		mockService.On("FindByID", swapID).Return(swap, nil)
		mockService.
			On("UpdateStatus", swapID, domain.StatusAccepted).
			Return(errors.New("DB error"))

		body := map[string]string{"status": "accepted"}
		jsonBody, _ := json.Marshal(body)

		req := httptest.NewRequest(http.MethodPatch, "/swap-requests/update-status/"+swapID.String(), bytes.NewBuffer(jsonBody))
		req.Header.Set("Content-Type", "application/json")
		resp := httptest.NewRecorder()

		router.ServeHTTP(resp, req)
		assert.Equal(t, http.StatusInternalServerError, resp.Code)
	})
}
