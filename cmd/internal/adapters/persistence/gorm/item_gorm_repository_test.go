package gorm_test

import (
	"github.com/google/uuid"
	"github.com/stretchr/testify/assert"
	"swapp-go/cmd/internal/adapters/persistence/gorm"
	"swapp-go/cmd/internal/adapters/persistence/gorm/testutils"
	"swapp-go/cmd/internal/adapters/persistence/models"
	"swapp-go/cmd/internal/domain"
	"testing"
)

func createTestItem(userID uuid.UUID) *domain.Item {
	return &domain.Item{
		Name:        "Test Item",
		Description: "A test item",
		PictureURL:  "/uploads/test.jpg",
		UserID:      userID,
	}
}

func TestCreateItem(t *testing.T) {
	db := testutils.SetupTestDB(t, &models.ItemModel{})
	repo := gorm.NewItemGormRepository(db)

	item := createTestItem(uuid.New())
	err := repo.CreateItem(item)

	assert.NoError(t, err)
	assert.NotEqual(t, uuid.Nil, item.ID)

	var dbItem models.ItemModel

	err = db.First(&dbItem, "id = ?", item.ID).Error
	assert.NoError(t, err)
	assert.Equal(t, item.Name, dbItem.Name)
}

func TestGetItemByID(t *testing.T) {
	db := testutils.SetupTestDB(t, &models.ItemModel{})
	repo := gorm.NewItemGormRepository(db)

	item := createTestItem(uuid.New())
	err := repo.CreateItem(item)
	assert.NoError(t, err)

	result, err := repo.GetItemByID(item.ID)
	assert.NoError(t, err)
	assert.Equal(t, item.Name, result.Name)
}

func TestUpdateItem(t *testing.T) {
	db := testutils.SetupTestDB(t, &models.ItemModel{})
	repo := gorm.NewItemGormRepository(db)

	item := createTestItem(uuid.New())
	err := repo.CreateItem(item)
	assert.NoError(t, err)

	updatedName := "Updated Name"
	fields := map[string]interface{}{"name": updatedName}

	updatedItem, err := repo.UpdateItem(item.ID, fields)
	assert.NoError(t, err)
	assert.Equal(t, updatedName, updatedItem.Name)
}

func TestDeleteItem(t *testing.T) {
	db := testutils.SetupTestDB(t, &models.ItemModel{})
	repo := gorm.NewItemGormRepository(db)

	item := createTestItem(uuid.New())
	err := repo.CreateItem(item)
	assert.NoError(t, err)

	err = repo.DeleteItem(item.ID)
	assert.NoError(t, err)

	_, err = repo.GetItemByID(item.ID)
	assert.Error(t, err)
}
