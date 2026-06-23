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
	"swapp-go/cmd/internal/adapters/handlers/mocks"
	"swapp-go/cmd/internal/domain"
	"swapp-go/cmd/internal/validators"
	"testing"
)

type mapStrStr map[string]string

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

// Test Helpers
func setupRouter(handler *handlers.UserHandler) *gin.Engine {
	gin.SetMode(gin.TestMode)

	validators.Init()

	router := gin.Default()
	router.POST("/users/register", handler.RegisterUser)
	router.POST("/users/login", handler.LoginUser)
	router.PATCH("/users/update", func(context *gin.Context) {
		context.Set("userID", uuid.Nil.String())
		handler.Update(context)
	})
	router.DELETE("/users/delete", func(context *gin.Context) {
		context.Set("userID", uuid.Nil.String())
		handler.Delete(context)
	})

	return router
}

func setupTest(t *testing.T) (*mocks.MockUserService, *gin.Engine) {
	t.Helper()

	mockService := new(mocks.MockUserService)
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

func TestUserHandler(t *testing.T) {
	t.Run("RegisterUser", func(t *testing.T) {
		t.Run("success", func(t *testing.T) {
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
		})

		t.Run("username_exists", func(t *testing.T) {
			mockService, router := setupTest(t)

			mockService.
				On("RegisterUser", mock.AnythingOfType("*domain.User")).
				Return(errors.New(expectedUsernameErr))

			response := performRequest(t, router, http.MethodPost, "/users/register", domainUser)
			assert.Equal(t, http.StatusBadRequest, response.Code)

			resp := parseErrorResponse(t, response)
			assert.Equal(t, invalidRequestErr, resp.Error)
			assert.Equal(t, expectedUsernameErr, resp.Details)

			mockService.AssertExpectations(t)
		})

		t.Run("email_exists", func(t *testing.T) {
			mockService, router := setupTest(t)

			mockService.
				On("RegisterUser", mock.AnythingOfType("*domain.User")).
				Return(errors.New(expectedEmailErr))

			response := performRequest(t, router, http.MethodPost, "/users/register", domainUser)
			assert.Equal(t, http.StatusBadRequest, response.Code)

			resp := parseErrorResponse(t, response)
			assert.Equal(t, invalidRequestErr, resp.Error)
			assert.Equal(t, expectedEmailErr, resp.Details)

			mockService.AssertExpectations(t)
		})

		t.Run("invalid_json", func(t *testing.T) {
			_, router := setupTest(t)

			response := performRequest(t, router, http.MethodPost, "/users/register", "{invalid-json")
			assert.Equal(t, http.StatusBadRequest, response.Code)

			resp := parseErrorResponse(t, response)
			assert.Equal(t, invalidRequestErr, resp.Error)
			assert.Contains(t, resp.Details, "invalid character")
		})
	})

	t.Run("LoginUser", func(t *testing.T) {
		t.Run("success", func(t *testing.T) {
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

			loginResp := parseLoginResponse(t, response)
			assert.Equal(t, testUser.ID.String(), loginResp.UserID)
			assert.Equal(t, testUser.Username, loginResp.Username)
			assert.Equal(t, testUser.Email, loginResp.Email)
			assert.Equal(t, testUser.Phone, loginResp.Phone)
			assert.Equal(t, testUser.Address, loginResp.Address)
			assert.Equal(t, testToken, loginResp.Token)

			mockService.AssertExpectations(t)
		})

		t.Run("invalid_credentials", func(t *testing.T) {
			mockService, router := setupTest(t)

			mockService.
				On("Authenticate", username, password).
				Return("", nil, errors.New("invalid credentials"))

			response := performRequest(t, router, http.MethodPost, "/users/login", loginPayload(username, password))
			assert.Equal(t, http.StatusUnauthorized, response.Code)

			errResp := parseErrorResponse(t, response)
			assert.Equal(t, "Unauthorized", errResp.Error)
			assert.Equal(t, "invalid credentials", errResp.Details)

			mockService.AssertExpectations(t)
		})

		t.Run("invalid_json", func(t *testing.T) {
			_, router := setupTest(t)

			response := performRequest(t, router, http.MethodPost, "/users/login", "{not-json")
			assert.Equal(t, http.StatusBadRequest, response.Code)

			errResp := parseErrorResponse(t, response)
			assert.Equal(t, invalidRequestErr, errResp.Error)
			assert.Contains(t, errResp.Details, "invalid character")
		})
	})

	t.Run("UpdateUser", func(t *testing.T) {
		t.Run("success", func(t *testing.T) {
			mockService, router := setupTest(t)

			mockService.
				On("Update", uuid.Nil, mock.MatchedBy(func(fields map[string]interface{}) bool {
					return fields["username"] == updatedUsername &&
						fields["email"] == updatedEmail &&
						fields["phone"] == updatedPhone &&
						fields["address"] == updatedAddress
				})).
				Return(updatedUser, nil)

			response := performRequest(t, router, http.MethodPatch, "/users/update", updatePayload)
			assert.Equal(t, http.StatusOK, response.Code)

			var parsed UserResponse
			err := json.Unmarshal(response.Body.Bytes(), &parsed)
			assert.NoError(t, err)
			assert.Equal(t, updateUserSuccessMsg, parsed.Message)

			mockService.AssertExpectations(t)
		})

		t.Run("failure", func(t *testing.T) {
			mockService := new(mocks.MockUserService)
			handler := handlers.NewUserHandler(mockService)

			router := gin.Default()
			router.PATCH("/users/update", func(c *gin.Context) {
				c.Set("userID", uuid.Nil.String())
				handler.Update(c)
			})

			mockService.
				On("Update", uuid.Nil, mock.Anything).
				Return(nil, errors.New("user not found"))

			body, _ := json.Marshal(updatePayload)
			req, _ := http.NewRequest(http.MethodPatch, "/users/update", bytes.NewReader(body))
			req.Header.Set("Content-Type", "application/json")

			resp := httptest.NewRecorder()
			router.ServeHTTP(resp, req)

			assert.Equal(t, http.StatusInternalServerError, resp.Code)

			var errResp ErrorResponse
			err := json.Unmarshal(resp.Body.Bytes(), &errResp)
			assert.NoError(t, err)
			assert.Equal(t, "Failed to update user", errResp.Error)
			assert.Equal(t, "user not found", errResp.Details)

			mockService.AssertExpectations(t)
		})

		t.Run("invalid_phone", func(t *testing.T) {
			_, router := setupTest(t)

			invalidPayload := map[string]interface{}{
				"username": updatedUsername,
				"email":    updatedEmail,
				"phone":    "1234567890",
				"address":  updatedAddress,
			}

			response := performRequest(t, router, http.MethodPatch, "/users/update", invalidPayload)
			assert.Equal(t, http.StatusBadRequest, response.Code)

			errResp := parseErrorResponse(t, response)
			assert.Equal(t, invalidRequestErr, errResp.Error)
			assert.Contains(t, errResp.Details, "Phone")
		})
	})

	t.Run("DeleteUser", func(t *testing.T) {
		t.Run("success", func(t *testing.T) {
			mockService, router := setupTest(t)

			mockService.
				On("Delete", uuid.Nil).
				Return(nil)

			response := performRequest(t, router, http.MethodDelete, "/users/delete", nil)
			assert.Equal(t, http.StatusOK, response.Code)

			var parsed map[string]string
			err := json.Unmarshal(response.Body.Bytes(), &parsed)
			assert.NoError(t, err)
			assert.Equal(t, "User deleted successfully!", parsed["message"])

			mockService.AssertExpectations(t)
		})

		t.Run("failure", func(t *testing.T) {
			mockService, router := setupTest(t)

			mockService.
				On("Delete", uuid.Nil).
				Return(errors.New("something went wrong"))

			response := performRequest(t, router, http.MethodDelete, "/users/delete", nil)
			assert.Equal(t, http.StatusInternalServerError, response.Code)

			errResp := parseErrorResponse(t, response)
			assert.Equal(t, "Failed to delete user", errResp.Error)
			assert.Equal(t, "something went wrong", errResp.Details)

			mockService.AssertExpectations(t)
		})
	})
}
