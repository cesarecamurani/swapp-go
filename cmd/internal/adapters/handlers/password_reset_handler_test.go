package handlers_test

import (
	"bytes"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"swapp-go/cmd/internal/adapters/handlers"
	"swapp-go/cmd/internal/adapters/persistence"
	"swapp-go/cmd/internal/application/service"
	"swapp-go/cmd/internal/domain"
	"swapp-go/cmd/internal/utils"
	"testing"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
	"github.com/stretchr/testify/assert"
	"gorm.io/driver/sqlite"
	"gorm.io/gorm"
)

func setupPasswordResetTestEnv(t *testing.T) (*gin.Engine, *gorm.DB, *domain.User, string) {
	gin.SetMode(gin.TestMode)
	db, err := gorm.Open(sqlite.Open(":memory:"), &gorm.Config{})
	assert.NoError(t, err)

	err = db.AutoMigrate(&persistence.UserModel{}, &persistence.PasswordResetModel{})
	if err != nil {
		return nil, nil, nil, ""
	}

	userRepo := persistence.NewGormUserRepository(db)
	resetRepo := persistence.NewGormPasswordResetRepository(db)
	userService := service.NewUserService(userRepo)
	resetService := service.NewPasswordResetService(resetRepo)

	handler := handlers.NewPasswordResetHandler(resetService, userService)
	router := gin.Default()
	router.POST("/request-reset", handler.RequestReset)
	router.POST("/reset-password", handler.ResetPassword)

	encryptedPassword, _ := utils.HashPassword("originalPassword")
	user := &domain.User{
		Username: "reset_user",
		Email:    "reset@example.com",
		Password: encryptedPassword,
	}
	err = userService.RegisterUser(user)
	assert.NoError(t, err)

	return router, db, user, "originalPassword"
}

func TestRequestReset_Success(t *testing.T) {
	router, _, user, _ := setupPasswordResetTestEnv(t)

	payload := map[string]string{"email": user.Email}
	jsonValue, _ := json.Marshal(payload)
	request, _ := http.NewRequest("POST", "/request-reset", bytes.NewBuffer(jsonValue))
	request.Header.Set("Content-Type", "application/json")

	resp := httptest.NewRecorder()
	router.ServeHTTP(resp, request)

	assert.Equal(t, http.StatusOK, resp.Code)
	assert.Contains(t, resp.Body.String(), "token")
}

func TestResetPassword_Success(t *testing.T) {
	router, db, user, _ := setupPasswordResetTestEnv(t)
	resetRepo := persistence.NewGormPasswordResetRepository(db)
	resetService := service.NewPasswordResetService(resetRepo)

	token, err := resetService.GenerateAndSaveToken(user.ID)
	assert.NoError(t, err)

	payload := map[string]string{
		"token":        token,
		"new_password": "newSecurePass123",
	}
	jsonValue, _ := json.Marshal(payload)
	request, _ := http.NewRequest("POST", "/reset-password", bytes.NewBuffer(jsonValue))
	request.Header.Set("Content-Type", "application/json")

	response := httptest.NewRecorder()
	router.ServeHTTP(response, request)

	assert.Equal(t, http.StatusOK, response.Code)
	assert.Contains(t, response.Body.String(), "Password reset successfully")
}

func TestResetPassword_InvalidToken(t *testing.T) {
	router, _, _, _ := setupPasswordResetTestEnv(t)
	payload := map[string]string{
		"token":        "invalid_token",
		"new_password": "NewPassword123",
	}
	jsonValue, _ := json.Marshal(payload)
	request, _ := http.NewRequest("POST", "/reset-password", bytes.NewBuffer(jsonValue))
	request.Header.Set("Content-Type", "application/json")

	response := httptest.NewRecorder()
	router.ServeHTTP(response, request)

	assert.Equal(t, http.StatusBadRequest, response.Code)
	assert.Contains(t, response.Body.String(), "Invalid or expired token")
}

func TestResetPassword_ExpiredToken(t *testing.T) {
	router, db, user, _ := setupPasswordResetTestEnv(t)
	resetRepo := persistence.NewGormPasswordResetRepository(db)

	expired := &domain.PasswordReset{
		Token:     uuid.NewString(),
		UserID:    user.ID,
		ExpiresAt: time.Now().Add(-1 * time.Hour),
	}
	err := resetRepo.Save(expired)
	assert.NoError(t, err)

	payload := map[string]string{
		"token":        expired.Token,
		"new_password": "NewPassword123",
	}
	jsonValue, _ := json.Marshal(payload)
	request, _ := http.NewRequest("POST", "/reset-password", bytes.NewBuffer(jsonValue))
	request.Header.Set("Content-Type", "application/json")

	response := httptest.NewRecorder()
	router.ServeHTTP(response, request)

	assert.Equal(t, http.StatusBadRequest, response.Code)
	assert.Contains(t, response.Body.String(), "Token expired")
}
