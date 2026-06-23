package gorm_test

import (
	"github.com/google/uuid"
	"github.com/stretchr/testify/assert"
	gormRepo "swapp-go/cmd/internal/adapters/persistence/gorm"
	"swapp-go/cmd/internal/adapters/persistence/gorm/testutils"
	"swapp-go/cmd/internal/adapters/persistence/models"
	"swapp-go/cmd/internal/domain"
	"testing"
	"time"
)

func TestPasswordResetRepository(t *testing.T) {
	db := testutils.SetupTestDB(t, &models.PasswordResetModel{})
	repo := gormRepo.NewPasswordResetGormRepository(db)

	t.Run("SaveAndGet", func(t *testing.T) {
		userID := uuid.New()
		token := "test_token"
		expiresAt := time.Now().Add(1 * time.Hour)

		reset := &domain.PasswordReset{
			Token:     token,
			UserID:    userID,
			ExpiresAt: expiresAt,
		}

		err := repo.Save(reset)
		assert.NoError(t, err)

		retrieved, err := repo.GetByToken(token)
		assert.NoError(t, err)
		assert.NotNil(t, retrieved)
		assert.Equal(t, reset.Token, retrieved.Token)
		assert.Equal(t, reset.UserID, retrieved.UserID)
		assert.WithinDuration(t, reset.ExpiresAt, retrieved.ExpiresAt, time.Second)
	})

	t.Run("GetByToken_NotFound", func(t *testing.T) {
		result, err := repo.GetByToken("non_existent_token")
		assert.Error(t, err)
		assert.Nil(t, result)
	})

	t.Run("Delete", func(t *testing.T) {
		userID := uuid.New()
		token := "some_token"

		reset := &domain.PasswordReset{
			Token:     token,
			UserID:    userID,
			ExpiresAt: time.Now().Add(1 * time.Hour),
		}

		err := repo.Save(reset)
		assert.NoError(t, err)

		err = repo.Delete(token)
		assert.NoError(t, err)

		_, err = repo.GetByToken(token)
		assert.Error(t, err)
	})
}
