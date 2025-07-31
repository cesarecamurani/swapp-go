package services

import (
	"github.com/google/uuid"
	"swapp-go/cmd/internal/application/ports"
	"swapp-go/cmd/internal/domain"
)

type ItemService struct {
	repo ports.ItemRepository
}

func NewItemService(repo ports.ItemRepository) *ItemService {
	return &ItemService{repo: repo}
}

func (itemService *ItemService) CreateItem(item *domain.Item) error {
	return itemService.repo.CreateItem(item)
}

func (itemService *ItemService) UpdateItem(id uuid.UUID, fields map[string]interface{}) (*domain.Item, error) {
	_, err := itemService.repo.GetItemByID(id)
	if err != nil {
		return nil, err
	}

	updatedItem, err := itemService.repo.UpdateItem(id, fields)
	if err != nil {
		return nil, err
	}

	return updatedItem, nil
}

func (itemService *ItemService) DeleteItem(id uuid.UUID) error {
	return itemService.repo.DeleteItem(id)
}

func (itemService *ItemService) GetItemByID(id uuid.UUID) (*domain.Item, error) {
	return itemService.repo.GetItemByID(id)
}
