package services_test

import (
	"errors"
	"github.com/google/uuid"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
	"swapp-go/cmd/internal/application/services"
	"swapp-go/cmd/internal/domain"
	"testing"
	"time"
)

var (
	validToken   = "valid_token"
	invalidToken = "invalid_token"
)

// Mocks
type MockPasswordResetRepository struct {
	mock.Mock
}

func (m *MockPasswordResetRepository) Save(token *domain.PasswordReset) error {
	args := m.Called(token)
	return args.Error(0)
}

func (m *MockPasswordResetRepository) GetByToken(token string) (*domain.PasswordReset, error) {
	args := m.Called(token)
	result := args.Get(0)
	if result == nil {
		return nil, args.Error(1)
	}
	return result.(*domain.PasswordReset), args.Error(1)
}

func (m *MockPasswordResetRepository) Delete(token string) error {
	args := m.Called(token)
	return args.Error(0)
}

// Tests
func TestGenerateAndSaveToken_Success(t *testing.T) {
	repo := new(MockPasswordResetRepository)
	resetService := services.NewPasswordResetService(repo)
	userID := uuid.New()

	repo.On("Save", mock.AnythingOfType("*domain.PasswordReset")).Return(nil)
	token, err := resetService.GenerateAndSaveToken(userID)

	assert.NoError(t, err)
	assert.NotEmpty(t, token)
	repo.AssertCalled(t, "Save", mock.AnythingOfType("*domain.PasswordReset"))
}

func TestValidateToken_Valid(t *testing.T) {
	repo := new(MockPasswordResetRepository)
	resetService := services.NewPasswordResetService(repo)

	expiredAt := time.Now().Add(1 * time.Hour)
	resetToken := &domain.PasswordReset{
		Token:     validToken,
		UserID:    uuid.New(),
		ExpiresAt: expiredAt,
	}
	repo.On("GetByToken", validToken).Return(resetToken, nil)

	token, err := resetService.ValidateToken(validToken)

	assert.NoError(t, err)
	assert.Equal(t, resetToken, token)
}

func TestValidateToken_Invalid(t *testing.T) {
	repo := new(MockPasswordResetRepository)
	resetService := services.NewPasswordResetService(repo)

	repo.On("GetByToken", invalidToken).Return(nil, errors.New("not found"))

	token, err := resetService.ValidateToken(invalidToken)

	assert.Error(t, err)
	assert.Nil(t, token)
}

func TestDeleteToken_Success(t *testing.T) {
	repo := new(MockPasswordResetRepository)
	resetService := services.NewPasswordResetService(repo)

	repo.On("Delete", validToken).Return(nil)

	err := resetService.DeleteToken(validToken)

	assert.NoError(t, err)
	repo.AssertCalled(t, "Delete", validToken)
}
