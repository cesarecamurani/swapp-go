package persistence_test

import (
	"github.com/google/uuid"
	"github.com/stretchr/testify/assert"
	"gorm.io/driver/sqlite"
	"gorm.io/gorm"
	"swapp-go/cmd/internal/adapters/persistence"
	"swapp-go/cmd/internal/domain"
	"testing"
	"time"
)

var validToken = "test_token"

func setupPasswordResetTestDB(t *testing.T) *gorm.DB {
	db, err := gorm.Open(sqlite.Open("file::memory:?cache=shared"), &gorm.Config{})
	assert.NoError(t, err)

	err = db.AutoMigrate(&persistence.PasswordResetModel{})
	assert.NoError(t, err)

	return db
}

func TestPasswordResetRepository_SaveAndGet(t *testing.T) {
	db := setupPasswordResetTestDB(t)
	repo := persistence.NewGormPasswordResetRepository(db)

	userID := uuid.New()
	token := validToken
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
}

func TestPasswordResetRepository_GetByToken_NotFound(t *testing.T) {
	db := setupPasswordResetTestDB(t)
	repo := persistence.NewGormPasswordResetRepository(db)

	result, err := repo.GetByToken("non_existent_token")
	assert.Error(t, err)
	assert.Nil(t, result)
}

func TestPasswordResetRepository_Delete(t *testing.T) {
	db := setupPasswordResetTestDB(t)
	repo := persistence.NewGormPasswordResetRepository(db)

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
}
