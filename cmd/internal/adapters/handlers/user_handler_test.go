package handlers_test

import (
	"bytes"
	"encoding/json"
	"errors"
	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
	"net/http"
	"net/http/httptest"
	"swapp-go/cmd/internal/adapters/handlers"
	"swapp-go/cmd/internal/application/service"
	"swapp-go/cmd/internal/domain"
	"testing"
)

type MockUserService struct {
	mock.Mock
}

type ErrorResponse struct {
	Error   string `json:"error"`
	Details string `json:"details,omitempty"`
}

var _ service.UserServiceInterface = (*MockUserService)(nil)
var username = "test_user"
var email = "test@example.com"
var password = "password123"
var domainUser = domain.User{
	Username: username,
	Email:    email,
	Password: password,
}

const (
	invalidRequestErr   = "Invalid request"
	expectedUsernameErr = "username already exists"
	expectedEmailErr    = "email already exists"
)

func (m *MockUserService) RegisterUser(user *domain.User) error {
	return m.Called(user).Error(0)
}

func (m *MockUserService) GetUserByID(id uuid.UUID) (*domain.User, error) {
	args := m.Called(id)

	if user, ok := args.Get(0).(*domain.User); ok {
		return user, args.Error(1)
	}

	return nil, args.Error(1)
}

func (m *MockUserService) GetUserByUsername(username string) (*domain.User, error) {
	args := m.Called(username)

	if user, ok := args.Get(0).(*domain.User); ok {
		return user, args.Error(1)
	}

	return nil, args.Error(1)
}

func (m *MockUserService) GetUserByEmail(email string) (*domain.User, error) {
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

func setupRouter(handler *handlers.UserHandler) *gin.Engine {
	gin.SetMode(gin.TestMode)

	router := gin.Default()
	router.POST("/users/register", handler.RegisterUser)

	return router
}

func setupTest(t *testing.T) (*MockUserService, *gin.Engine) {
	t.Helper()

	mockService := new(MockUserService)
	handler := handlers.NewUserHandler(mockService)
	router := setupRouter(handler)

	return mockService, router
}

func performPostRequest(t *testing.T, router *gin.Engine, url string, body interface{}) *httptest.ResponseRecorder {
	t.Helper()

	jsonBody, err := json.Marshal(body)
	assert.NoError(t, err)

	request, err := http.NewRequest(http.MethodPost, url, bytes.NewBuffer(jsonBody))
	assert.NoError(t, err)

	request.Header.Set("Content-Type", "application/json")

	response := httptest.NewRecorder()
	router.ServeHTTP(response, request)

	return response
}

func parseErrorResponse(t *testing.T, response *httptest.ResponseRecorder) *ErrorResponse {
	t.Helper()

	var errResponse ErrorResponse
	err := json.Unmarshal(response.Body.Bytes(), &errResponse)
	assert.NoError(t, err)

	return &errResponse
}

func TestRegisterUser_Success(t *testing.T) {
	mockService, router := setupTest(t)

	mockService.
		On("RegisterUser", mock.AnythingOfType("*domain.User")).
		Return(nil)

	userPayload := map[string]string{
		"username": username,
		"email":    email,
		"password": password,
	}

	response := performPostRequest(t, router, "/users/register", userPayload)

	assert.Equal(t, http.StatusCreated, response.Code)
	mockService.AssertExpectations(t)
}

func TestRegisterUser_UsernameExists(t *testing.T) {
	mockService, router := setupTest(t)

	mockService.
		On("RegisterUser", mock.AnythingOfType("*domain.User")).
		Return(errors.New(expectedUsernameErr))

	response := performPostRequest(t, router, "/users/register", domainUser)
	assert.Equal(t, http.StatusBadRequest, response.Code)

	parsedResponse := parseErrorResponse(t, response)
	assert.Equal(t, invalidRequestErr, parsedResponse.Error)
	assert.Equal(t, expectedUsernameErr, parsedResponse.Details)

	mockService.AssertExpectations(t)
}

func TestRegisterUser_EmailAlreadyExists(t *testing.T) {
	mockService, router := setupTest(t)

	mockService.
		On("RegisterUser", mock.AnythingOfType("*domain.User")).
		Return(errors.New(expectedEmailErr))

	response := performPostRequest(t, router, "/users/register", domainUser)
	assert.Equal(t, http.StatusBadRequest, response.Code)

	parsedResponse := parseErrorResponse(t, response)
	assert.Equal(t, invalidRequestErr, parsedResponse.Error)
	assert.Equal(t, expectedEmailErr, parsedResponse.Details)

	mockService.AssertExpectations(t)
}

func TestRegisterUser_InvalidJSON(t *testing.T) {
	_, router := setupTest(t)

	request, _ := http.NewRequest(http.MethodPost, "/users/register", bytes.NewBufferString("{invalid-json"))
	request.Header.Set("Content-Type", "application/json")
	response := httptest.NewRecorder()
	router.ServeHTTP(response, request)

	assert.Equal(t, http.StatusBadRequest, response.Code)

	parsedResponse := parseErrorResponse(t, response)
	assert.Equal(t, invalidRequestErr, parsedResponse.Error)
	assert.Contains(t, parsedResponse.Details, "invalid character")
}
