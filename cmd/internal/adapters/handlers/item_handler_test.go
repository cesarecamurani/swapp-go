package handlers_test

import (
	"bytes"
	"errors"
	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
	"io"
	"mime/multipart"
	"net/http"
	"net/http/httptest"
	"swapp-go/cmd/internal/adapters/handlers"
	"swapp-go/cmd/internal/domain"
	"testing"
)

// Mocks
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

// Test Helpers
func writeFile(t *testing.T, w io.Writer, data []byte) {
	t.Helper()

	if _, err := w.Write(data); err != nil {
		t.Fatalf("write failed: %v", err)
	}
}

func writeFormField(t *testing.T, writer *multipart.Writer, fieldName, value string) {
	t.Helper()

	if err := writer.WriteField(fieldName, value); err != nil {
		t.Fatalf("failed to write form field %q: %v", fieldName, err)
	}
}

func closeWriter(t *testing.T, c io.Closer) {
	t.Helper()

	if err := c.Close(); err != nil {
		t.Fatalf("close failed: %v", err)
	}
}

// Tests
func TestCreateItem_Success(t *testing.T) {
	gin.SetMode(gin.TestMode)

	mockService := new(MockItemService)
	handler := handlers.NewItemHandler(mockService)

	userID := uuid.New()

	mockService.On("Create", mock.AnythingOfType("*domain.Item")).Return(nil)

	bodyBuffer := &bytes.Buffer{}
	formWriter := multipart.NewWriter(bodyBuffer)
	_ = formWriter.WriteField("name", "Test Item")
	_ = formWriter.WriteField("description", "A good item")

	fileWriter, _ := formWriter.CreateFormFile("picture", "image.jpg")
	writeFile(t, fileWriter, []byte("fake image content"))
	closeWriter(t, formWriter)

	request := httptest.NewRequest(http.MethodPost, "/items/create", bodyBuffer)
	request.Header.Set("Content-Type", formWriter.FormDataContentType())

	responseRecorder := httptest.NewRecorder()
	context, _ := gin.CreateTestContext(responseRecorder)
	context.Request = request
	context.Set("userID", userID.String())

	handler.Create(context)

	assert.Equal(t, http.StatusCreated, responseRecorder.Code)
	assert.Contains(t, responseRecorder.Body.String(), "Item created successfully!")
	mockService.AssertCalled(t, "Create", mock.AnythingOfType("*domain.Item"))
}

func TestCreateItem_InvalidUserID(t *testing.T) {
	mockService := new(MockItemService)
	handler := handlers.NewItemHandler(mockService)

	bodyBuffer := &bytes.Buffer{}
	formWriter := multipart.NewWriter(bodyBuffer)
	writeFormField(t, formWriter, "name", "Some Item")
	writeFormField(t, formWriter, "description", "Description")
	closeWriter(t, formWriter)

	request := httptest.NewRequest(http.MethodPost, "/items/create", bodyBuffer)
	request.Header.Set("Content-Type", formWriter.FormDataContentType())

	responseRecorder := httptest.NewRecorder()
	context, _ := gin.CreateTestContext(responseRecorder)
	context.Request = request
	context.Set("userID", "invalid-uuid")

	handler.Create(context)

	assert.Equal(t, http.StatusBadRequest, responseRecorder.Code)
	assert.Contains(t, responseRecorder.Body.String(), "Invalid user ID")
}

func TestUpdateItem_Success(t *testing.T) {
	mockService := new(MockItemService)
	handler := handlers.NewItemHandler(mockService)

	itemID := uuid.New()
	userID := uuid.New()

	existingItem := &domain.Item{
		ID:     itemID,
		UserID: userID,
	}

	updatedItem := &domain.Item{
		ID:          itemID,
		Name:        "Updated Name",
		Description: "Updated Description",
		PictureURL:  "/uploads/new.jpg",
		UserID:      userID,
	}

	mockService.On("FindByID", itemID).Return(existingItem, nil)
	mockService.On("Update", itemID, mock.AnythingOfType("map[string]interface {}")).Return(updatedItem, nil)

	bodyBuffer := &bytes.Buffer{}
	formWriter := multipart.NewWriter(bodyBuffer)
	_ = formWriter.WriteField("name", "Updated Name")
	_ = formWriter.WriteField("description", "Updated Description")

	fileWriter, _ := formWriter.CreateFormFile("picture", "new.jpg")
	writeFile(t, fileWriter, []byte("new image content"))
	closeWriter(t, formWriter)

	request := httptest.NewRequest(http.MethodPost, "/items/"+itemID.String(), bodyBuffer)
	request.Header.Set("Content-Type", formWriter.FormDataContentType())

	responseRecorder := httptest.NewRecorder()
	context, _ := gin.CreateTestContext(responseRecorder)
	context.Request = request
	context.Params = gin.Params{{Key: "id", Value: itemID.String()}}
	context.Set("userID", userID.String())

	handler.Update(context)

	assert.Equal(t, http.StatusOK, responseRecorder.Code)
	assert.Contains(t, responseRecorder.Body.String(), "Item updated successfully")
}

func TestDeleteItem_Success(t *testing.T) {
	mockService := new(MockItemService)
	handler := handlers.NewItemHandler(mockService)

	itemID := uuid.New()
	userID := uuid.New()

	item := &domain.Item{
		ID:     itemID,
		UserID: userID,
	}

	mockService.On("FindByID", itemID).Return(item, nil)
	mockService.On("Delete", itemID).Return(nil)

	request := httptest.NewRequest(http.MethodDelete, "/items/"+itemID.String(), nil)
	responseRecorder := httptest.NewRecorder()
	context, _ := gin.CreateTestContext(responseRecorder)
	context.Request = request
	context.Params = gin.Params{{Key: "id", Value: itemID.String()}}
	context.Set("userID", userID.String())

	handler.Delete(context)

	assert.Equal(t, http.StatusOK, responseRecorder.Code)
	assert.Contains(t, responseRecorder.Body.String(), "Item deleted successfully!")
}

func TestDeleteItem_Unauthorized(t *testing.T) {
	mockService := new(MockItemService)
	handler := handlers.NewItemHandler(mockService)

	itemID := uuid.New()
	itemOwnerID := uuid.New()
	requestUserID := uuid.New()

	item := &domain.Item{
		ID:     itemID,
		UserID: itemOwnerID,
	}

	mockService.On("FindByID", itemID).Return(item, nil)

	request := httptest.NewRequest(http.MethodDelete, "/items/"+itemID.String(), nil)
	responseRecorder := httptest.NewRecorder()
	context, _ := gin.CreateTestContext(responseRecorder)
	context.Request = request
	context.Params = gin.Params{{Key: "id", Value: itemID.String()}}
	context.Set("userID", requestUserID.String())

	handler.Delete(context)

	assert.Equal(t, http.StatusUnauthorized, responseRecorder.Code)
	assert.Contains(t, responseRecorder.Body.String(), "doesn't belong to you")
}

func TestGetItemByID_Success(t *testing.T) {
	mockService := new(MockItemService)
	handler := handlers.NewItemHandler(mockService)

	itemID := uuid.New()
	userID := uuid.New()

	mockItem := &domain.Item{
		ID:          itemID,
		Name:        "Test",
		Description: "A test item",
		PictureURL:  "/uploads/test.jpg",
		UserID:      userID,
	}
	mockService.On("FindByID", itemID).Return(mockItem, nil)

	request := httptest.NewRequest(http.MethodGet, "/items/"+itemID.String(), nil)
	responseRecorder := httptest.NewRecorder()
	context, _ := gin.CreateTestContext(responseRecorder)
	context.Request = request
	context.Params = gin.Params{{Key: "id", Value: itemID.String()}}

	handler.FindByID(context)

	assert.Equal(t, http.StatusOK, responseRecorder.Code)
	assert.Contains(t, responseRecorder.Body.String(), itemID.String())
}

func TestGetItemByID_NotFound(t *testing.T) {
	mockService := new(MockItemService)
	handler := handlers.NewItemHandler(mockService)

	itemID := uuid.New()
	mockService.On("FindByID", itemID).Return(nil, errors.New("not found"))

	request := httptest.NewRequest(http.MethodGet, "/items/"+itemID.String(), nil)
	responseRecorder := httptest.NewRecorder()
	context, _ := gin.CreateTestContext(responseRecorder)
	context.Request = request
	context.Params = gin.Params{{Key: "id", Value: itemID.String()}}

	handler.FindByID(context)

	assert.Equal(t, http.StatusNotFound, responseRecorder.Code)
}
