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
	"swapp-go/cmd/internal/domain"
	"swapp-go/cmd/internal/validators"
	"testing"
)

type mapStrStr map[string]string

type MockUserService struct {
	mock.Mock
}

type ErrorResponse struct {
	Error   string `json:"error"`
	Details string `json:"details,omitempty"`
}

type UserResponse struct {
	Message string `json:"message"`
	User    struct {
		UserID   string `json:"user_id"`
		Username string `json:"username"`
		Email    string `json:"email"`
	} `json:"user"`
}

// Constants and shared variables
const (
	invalidRequestErr    = "Invalid request"
	expectedUsernameErr  = "username already exists"
	expectedEmailErr     = "email already exists"
	updateUserSuccessMsg = "User updated successfully!"
)

var (
	username  = "test_user"
	email     = "test@example.com"
	password  = "password123"
	phone     = "+447712345678"
	address   = "1, Main Street"
	testToken = "test-token"

	domainUser = domain.User{
		Username: username,
		Email:    email,
		Password: password,
	}

	updatedUsername = "updated_user"
	updatedEmail    = "updated@example.com"
	updatedPhone    = "+447787654321"
	updatedAddress  = "2, Main Street"

	updatePayload = map[string]interface{}{
		"username": updatedUsername,
		"email":    updatedEmail,
		"phone":    updatedPhone,
		"address":  updatedAddress,
	}

	updatedUser = &domain.User{
		ID:       uuid.Nil,
		Username: updatedUsername,
		Email:    updatedEmail,
		Phone:    &updatedPhone,
		Address:  &updatedAddress,
	}
)

// MockUserService methods
func (m *MockUserService) RegisterUser(user *domain.User) error {
	return m.Called(user).Error(0)
}

func (m *MockUserService) UpdateUser(id uuid.UUID, fields map[string]interface{}) (*domain.User, error) {
	args := m.Called(id, fields)

	if user, ok := args.Get(0).(*domain.User); ok {
		return user, args.Error(1)
	}

	return nil, args.Error(1)
}

func (m *MockUserService) DeleteUser(id uuid.UUID) error {
	return m.Called(id).Error(0)
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

// Helper Functions
func setupRouter(handler *handlers.UserHandler) *gin.Engine {
	gin.SetMode(gin.TestMode)

	validators.Init()

	router := gin.Default()
	router.POST("/users/register", handler.RegisterUser)
	router.POST("/users/login", handler.LoginUser)
	router.PATCH("/users/update", func(context *gin.Context) {
		context.Set("userID", uuid.Nil.String())
		handler.UpdateUser(context)
	})
	router.DELETE("/users/delete", func(context *gin.Context) {
		context.Set("userID", uuid.Nil.String())
		handler.DeleteUser(context)
	})

	return router
}

func setupTest(t *testing.T) (*MockUserService, *gin.Engine) {
	t.Helper()

	mockService := new(MockUserService)
	handler := handlers.NewUserHandler(mockService)
	router := setupRouter(handler)

	return mockService, router
}

func performRequest(t *testing.T, router *gin.Engine, method, url string, body interface{}) *httptest.ResponseRecorder {
	t.Helper()

	var jsonBody []byte
	var err error

	switch v := body.(type) {
	case string:
		jsonBody = []byte(v)
	case []byte:
		jsonBody = v
	default:
		jsonBody, err = json.Marshal(v)
		assert.NoError(t, err)
	}

	request, err := http.NewRequest(method, url, bytes.NewBuffer(jsonBody))
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

func loginPayload(username, password string) mapStrStr {
	return mapStrStr{
		"username": username,
		"password": password,
	}
}

func parseLoginResponse(t *testing.T, response *httptest.ResponseRecorder) *handlers.LoginUserResponse {
	t.Helper()

	var loginResponse handlers.LoginUserResponse
	err := json.Unmarshal(response.Body.Bytes(), &loginResponse)
	assert.NoError(t, err)

	return &loginResponse
}

// Tests
// RegisterUser
func TestRegisterUser_Success(t *testing.T) {
	mockService, router := setupTest(t)

	mockService.
		On("RegisterUser", mock.AnythingOfType("*domain.User")).
		Return(nil)

	userPayload := map[string]string{
		"username": username,
		"email":    email,
		"password": password,
		"phone":    phone,
		"address":  address,
	}

	response := performRequest(t, router, http.MethodPost, "/users/register", userPayload)

	assert.Equal(t, http.StatusCreated, response.Code)
	mockService.AssertExpectations(t)
}

func TestRegisterUser_UsernameExists(t *testing.T) {
	mockService, router := setupTest(t)

	mockService.
		On("RegisterUser", mock.AnythingOfType("*domain.User")).
		Return(errors.New(expectedUsernameErr))

	response := performRequest(t, router, http.MethodPost, "/users/register", domainUser)
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

	response := performRequest(t, router, http.MethodPost, "/users/register", domainUser)
	assert.Equal(t, http.StatusBadRequest, response.Code)

	parsedResponse := parseErrorResponse(t, response)
	assert.Equal(t, invalidRequestErr, parsedResponse.Error)
	assert.Equal(t, expectedEmailErr, parsedResponse.Details)

	mockService.AssertExpectations(t)
}

func TestRegisterUser_InvalidJSON(t *testing.T) {
	_, router := setupTest(t)

	response := performRequest(t, router, http.MethodPost, "/users/register", "{invalid-json")
	assert.Equal(t, http.StatusBadRequest, response.Code)

	parsedResponse := parseErrorResponse(t, response)
	assert.Equal(t, invalidRequestErr, parsedResponse.Error)
	assert.Contains(t, parsedResponse.Details, "invalid character")
}

// LoginUser
func TestLoginUser_Success(t *testing.T) {
	mockService, router := setupTest(t)

	userID := uuid.New()
	testUser := &domain.User{
		ID:       userID,
		Username: username,
		Email:    email,
		Password: password,
		Phone:    &phone,
		Address:  &address,
	}

	mockService.
		On("Authenticate", username, password).
		Return(testToken, testUser, nil)

	response := performRequest(t, router, http.MethodPost, "/users/login", loginPayload(username, password))
	assert.Equal(t, http.StatusOK, response.Code)

	loginResponse := parseLoginResponse(t, response)
	assert.Equal(t, testUser.ID.String(), loginResponse.UserID)
	assert.Equal(t, testUser.Username, loginResponse.Username)
	assert.Equal(t, testUser.Email, loginResponse.Email)
	assert.Equal(t, testUser.Phone, loginResponse.Phone)
	assert.Equal(t, testUser.Address, loginResponse.Address)
	assert.Equal(t, testToken, loginResponse.Token)

	mockService.AssertExpectations(t)
}

func TestLoginUser_InvalidCredentials(t *testing.T) {
	mockService, router := setupTest(t)

	mockService.
		On("Authenticate", username, password).
		Return("", nil, errors.New("invalid credentials"))

	response := performRequest(t, router, http.MethodPost, "/users/login", loginPayload(username, password))
	assert.Equal(t, http.StatusUnauthorized, response.Code)

	errResponse := parseErrorResponse(t, response)
	assert.Equal(t, "Unauthorized", errResponse.Error)
	assert.Equal(t, "invalid credentials", errResponse.Details)

	mockService.AssertExpectations(t)
}

func TestLoginUser_InvalidJSON(t *testing.T) {
	_, router := setupTest(t)

	response := performRequest(t, router, http.MethodPost, "/users/login", "{not-json")
	assert.Equal(t, http.StatusBadRequest, response.Code)

	errResp := parseErrorResponse(t, response)
	assert.Equal(t, invalidRequestErr, errResp.Error)
	assert.Contains(t, errResp.Details, "invalid character")
}

// UpdateUser
func TestUpdateUser_Success(t *testing.T) {
	mockService, router := setupTest(t)

	mockService.
		On("UpdateUser", uuid.Nil, mock.MatchedBy(func(fields map[string]interface{}) bool {
			return fields["username"] == updatedUsername &&
				fields["email"] == updatedEmail &&
				fields["phone"] == updatedPhone &&
				fields["address"] == updatedAddress
		})).
		Return(updatedUser, nil)

	response := performRequest(t, router, http.MethodPatch, "/users/update", updatePayload)

	assert.Equal(t, http.StatusOK, response.Code)

	var parsedResponse UserResponse

	err := json.Unmarshal(response.Body.Bytes(), &parsedResponse)

	assert.NoError(t, err)
	assert.Equal(t, updateUserSuccessMsg, parsedResponse.Message)

	mockService.AssertExpectations(t)
}

func TestUpdateUser_Failure(t *testing.T) {
	mockService := new(MockUserService)
	handler := handlers.NewUserHandler(mockService)

	router := gin.Default()
	router.PATCH("/users/update", func(context *gin.Context) {
		context.Set("userID", uuid.Nil.String())
		handler.UpdateUser(context)
	})

	expectedErr := errors.New("user not found")

	mockService.
		On("UpdateUser", uuid.Nil, mock.MatchedBy(func(fields map[string]interface{}) bool {
			return fields["username"] == updatedUsername &&
				fields["email"] == updatedEmail &&
				fields["phone"] == updatedPhone &&
				fields["address"] == updatedAddress
		})).
		Return(nil, expectedErr)

	body, _ := json.Marshal(updatePayload)
	request, _ := http.NewRequest(http.MethodPatch, "/users/update", bytes.NewReader(body))
	request.Header.Set("Content-Type", "application/json")

	response := httptest.NewRecorder()
	router.ServeHTTP(response, request)

	assert.Equal(t, http.StatusInternalServerError, response.Code)

	var errResponse ErrorResponse

	err := json.Unmarshal(response.Body.Bytes(), &errResponse)

	assert.NoError(t, err)
	assert.Equal(t, "Failed to update user", errResponse.Error)
	assert.Equal(t, "user not found", errResponse.Details)

	mockService.AssertExpectations(t)
}

func TestUpdateUser_InvalidPhone(t *testing.T) {
	_, router := setupTest(t)

	invalidPhonePayload := map[string]interface{}{
		"username": updatedUsername,
		"email":    updatedEmail,
		"phone":    "1234567890",
		"address":  updatedAddress,
	}

	response := performRequest(t, router, http.MethodPatch, "/users/update", invalidPhonePayload)

	assert.Equal(t, http.StatusBadRequest, response.Code)

	errResp := parseErrorResponse(t, response)
	assert.Equal(t, invalidRequestErr, errResp.Error)
	assert.Contains(t, errResp.Details, "Phone")
}

// DeleteUser
func TestDeleteUser_Success(t *testing.T) {
	mockService, router := setupTest(t)

	mockService.
		On("DeleteUser", uuid.Nil).
		Return(nil)

	response := performRequest(t, router, http.MethodDelete, "/users/delete", nil)

	assert.Equal(t, http.StatusOK, response.Code)

	var parsed map[string]string

	err := json.Unmarshal(response.Body.Bytes(), &parsed)

	assert.NoError(t, err)
	assert.Equal(t, "User deleted successfully!", parsed["message"])

	mockService.AssertExpectations(t)
}

func TestDeleteUser_Failure(t *testing.T) {
	mockService, router := setupTest(t)

	mockService.
		On("DeleteUser", uuid.Nil).
		Return(errors.New("something went wrong"))

	response := performRequest(t, router, http.MethodDelete, "/users/delete", nil)

	assert.Equal(t, http.StatusInternalServerError, response.Code)

	errResponse := parseErrorResponse(t, response)

	assert.Equal(t, "Failed to delete user", errResponse.Error)
	assert.Equal(t, "something went wrong", errResponse.Details)

	mockService.AssertExpectations(t)
}
