package services_test

import (
	"errors"
	"testing"
	"time"

	"github.com/google/uuid"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
	"swapp-go/cmd/internal/application/services"
	"swapp-go/cmd/internal/application/services/mocks"
	"swapp-go/cmd/internal/domain"
)

var (
	validToken   = "valid_token"
	invalidToken = "invalid_token"
)

func TestPasswordResetService(t *testing.T) {
	t.Run("GenerateAndSaveToken_Success", func(t *testing.T) {
		repo := new(mocks.MockPasswordResetRepository)
		resetService := services.NewPasswordResetService(repo)
		userID := uuid.New()

		repo.On("Save", mock.AnythingOfType("*domain.PasswordReset")).Return(nil)

		token, err := resetService.GenerateAndSaveToken(userID)

		assert.NoError(t, err)
		assert.NotEmpty(t, token)
		repo.AssertCalled(t, "Save", mock.AnythingOfType("*domain.PasswordReset"))
	})

	t.Run("ValidateToken_Valid", func(t *testing.T) {
		repo := new(mocks.MockPasswordResetRepository)
		resetService := services.NewPasswordResetService(repo)

		expiresAt := time.Now().Add(1 * time.Hour)
		resetToken := &domain.PasswordReset{
			Token:     validToken,
			UserID:    uuid.New(),
			ExpiresAt: expiresAt,
		}

		repo.On("GetByToken", validToken).Return(resetToken, nil)

		token, err := resetService.ValidateToken(validToken)

		assert.NoError(t, err)
		assert.Equal(t, resetToken, token)
	})

	t.Run("ValidateToken_Invalid", func(t *testing.T) {
		repo := new(mocks.MockPasswordResetRepository)
		resetService := services.NewPasswordResetService(repo)

		repo.On("GetByToken", invalidToken).Return(nil, errors.New("not found"))

		token, err := resetService.ValidateToken(invalidToken)

		assert.Error(t, err)
		assert.Nil(t, token)
	})

	t.Run("DeleteToken_Success", func(t *testing.T) {
		repo := new(mocks.MockPasswordResetRepository)
		resetService := services.NewPasswordResetService(repo)

		repo.On("Delete", validToken).Return(nil)

		err := resetService.DeleteToken(validToken)

		assert.NoError(t, err)
		repo.AssertCalled(t, "Delete", validToken)
	})
}
