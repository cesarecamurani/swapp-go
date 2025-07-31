package gorm

import (
	"github.com/google/uuid"
	"gorm.io/gorm"
	"swapp-go/cmd/internal/adapters/persistence/models"
	"swapp-go/cmd/internal/domain"
)

type ItemGormRepository struct {
	db *gorm.DB
}

func NewItemGormRepository(db *gorm.DB) *ItemGormRepository {
	return &ItemGormRepository{db}
}

func toItemModel(item *domain.Item) *models.ItemModel {
	id := item.ID
	if id == uuid.Nil {
		id = uuid.New()
	}

	return &models.ItemModel{
		ID:          id,
		Name:        item.Name,
		Description: item.Description,
		PictureURL:  item.PictureURL,
		UserID:      item.UserID,
	}
}

func toDomainItem(model *models.ItemModel) *domain.Item {
	return &domain.Item{
		ID:          model.ID,
		Name:        model.Name,
		Description: model.Description,
		PictureURL:  model.PictureURL,
		UserID:      model.UserID,
	}
}

func (itemGorm *ItemGormRepository) CreateItem(item *domain.Item) error {
	model := toItemModel(item)

	if result := itemGorm.db.Create(model); result.Error != nil {
		return result.Error
	}

	item.ID = model.ID

	return nil
}

func (itemGorm *ItemGormRepository) UpdateItem(id uuid.UUID, fields map[string]interface{}) (*domain.Item, error) {
	if err := itemGorm.db.Model(&models.ItemModel{}).Where("id = ?", id).Updates(fields).Error; err != nil {
		return nil, err
	}

	var updatedItemModel models.ItemModel
	if err := itemGorm.db.Where("id = ?", id).First(&updatedItemModel).Error; err != nil {
		return nil, err
	}

	return toDomainItem(&updatedItemModel), nil
}

func (itemGorm *ItemGormRepository) DeleteItem(id uuid.UUID) error {
	return itemGorm.db.Delete(&models.ItemModel{}, id).Error
}

func (itemGorm *ItemGormRepository) GetItemByID(id uuid.UUID) (*domain.Item, error) {
	var itemModel models.ItemModel

	if err := itemGorm.db.First(&itemModel, id).Error; err != nil {
		return nil, err
	}

	return toDomainItem(&itemModel), nil
}
