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

var (
	phone   = "+441234567890"
	address = "1, Main Street"
	user    = &domain.User{
		Username: "test_user",
		Password: "hashed_password",
		Email:    "test@email.com",
		Phone:    &phone,
		Address:  &address,
	}
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

	err := repo.CreateUser(user)
	assert.NoError(t, err)
	assert.NotEqual(t, uuid.Nil, user.ID)

	userByID, err := repo.GetUserByID(user.ID)
	assert.NoError(t, err)
	assert.Equal(t, user.Username, userByID.Username)
	assert.Equal(t, user.Email, userByID.Email)
	assert.NotNil(t, userByID.Phone)
	assert.NotNil(t, userByID.Address)
	assert.Equal(t, *user.Phone, *userByID.Phone)
	assert.Equal(t, *user.Address, *userByID.Address)

	userByUsername, err := repo.GetUserByUsername(user.Username)
	assert.NoError(t, err)
	assert.Equal(t, user.ID, userByUsername.ID)

	userByEmail, err := repo.GetUserByEmail(user.Email)
	assert.NoError(t, err)
	assert.Equal(t, user.ID, userByEmail.ID)
}

func TestGormUserRepository_UpdateUser(t *testing.T) {
	db := setupTestDB(t)
	repo := persistence.NewGormUserRepository(db)

	err := repo.CreateUser(user)
	assert.NoError(t, err)

	updatedPhone := "+44778654321"
	updatedAddress := "2, Main Street"
	updatedFields := map[string]interface{}{
		"username": "updated_user",
		"email":    "updated@example.com",
		"phone":    updatedPhone,
		"address":  updatedAddress,
	}

	updatedUser, err := repo.UpdateUser(user.ID, updatedFields)
	assert.NoError(t, err)
	assert.Equal(t, "updated_user", updatedUser.Username)
	assert.Equal(t, "updated@example.com", updatedUser.Email)
	assert.NotNil(t, updatedUser.Phone)
	assert.NotNil(t, updatedUser.Address)
	assert.Equal(t, updatedPhone, *updatedUser.Phone)
	assert.Equal(t, updatedAddress, *updatedUser.Address)
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
