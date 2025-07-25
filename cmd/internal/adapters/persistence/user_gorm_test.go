package persistence_test

import (
	"github.com/google/uuid"
	"github.com/stretchr/testify/assert"
	"gorm.io/driver/sqlite"
	"gorm.io/gorm"
	"swapp-go/cmd/internal/adapters/persistence"
	"swapp-go/cmd/internal/domain"
	"testing"
)

func setupTestDB(t *testing.T) *gorm.DB {
	db, err := gorm.Open(sqlite.Open(":memory:"), &gorm.Config{})
	assert.NoError(t, err)

	err = db.AutoMigrate(&persistence.UserModel{})
	assert.NoError(t, err)

	return db
}

func TestGormUserRepository_CreateAndGetUser(t *testing.T) {
	db := setupTestDB(t)
	repo := persistence.NewGormUserRepository(db)

	user := &domain.User{
		Username: "test_user",
		Password: "hashed_password",
		Email:    "test@email.com",
	}

	err := repo.CreateUser(user)
	assert.NoError(t, err)
	assert.NotEqual(t, uuid.Nil, user.ID)

	userByID, err := repo.GetUserByID(user.ID)
	assert.NoError(t, err)
	assert.Equal(t, user.Username, userByID.Username)
	assert.Equal(t, user.Email, userByID.Email)

	userByUsername, err := repo.GetUserByUsername(user.Username)
	assert.NoError(t, err)
	assert.Equal(t, user.ID, userByUsername.ID)

	userByEmail, err := repo.GetUserByEmail(user.Email)
	assert.NoError(t, err)
	assert.Equal(t, user.ID, userByEmail.ID)
}

func TestGormUserRepository_NotFound(t *testing.T) {
	db := setupTestDB(t)
	repo := persistence.NewGormUserRepository(db)

	randomID := uuid.New()
	user, err := repo.GetUserByID(randomID)
	assert.Error(t, err)
	assert.Nil(t, user)

	user, err = repo.GetUserByUsername("non_existent_username")
	assert.Error(t, err)
	assert.Nil(t, user)

	user, err = repo.GetUserByEmail("nonexistent.email@example.com")
	assert.Error(t, err)
	assert.Nil(t, user)
}
