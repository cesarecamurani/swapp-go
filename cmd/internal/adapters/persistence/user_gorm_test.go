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

var user = &domain.User{
	Username: "test_user",
	Password: "hashed_password",
	Email:    "test@email.com",
}

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
	notFoundUser, err := repo.GetUserByID(randomID)
	assert.Error(t, err)
	assert.Nil(t, notFoundUser)

	notFoundUser, err = repo.GetUserByUsername("non_existent_username")
	assert.Error(t, err)
	assert.Nil(t, notFoundUser)

	notFoundUser, err = repo.GetUserByEmail("nonexistent.email@example.com")
	assert.Error(t, err)
	assert.Nil(t, notFoundUser)
}

func TestDeleteUser(t *testing.T) {
	db := setupTestDB(t) // e.g. SQLite in-memory
	repo := persistence.NewGormUserRepository(db)

	assert.NoError(t, db.Create(user).Error)

	err := repo.DeleteUser(user.ID)
	assert.NoError(t, err)

	var found persistence.UserModel

	err = db.First(&found, "id = ?", user.ID).Error

	assert.ErrorIs(t, err, gorm.ErrRecordNotFound)
}
