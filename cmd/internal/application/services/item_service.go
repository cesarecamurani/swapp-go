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

func (itemService *ItemService) Create(item *domain.Item) error {
	return itemService.repo.Create(item)
}

func (itemService *ItemService) Update(id uuid.UUID, fields map[string]interface{}) (*domain.Item, error) {
	_, err := itemService.repo.FindByID(id)
	if err != nil {
		return nil, err
	}

	updatedItem, err := itemService.repo.Update(id, fields)
	if err != nil {
		return nil, err
	}

	return updatedItem, nil
}

func (itemService *ItemService) Delete(id uuid.UUID) error {
	return itemService.repo.Delete(id)
}

func (itemService *ItemService) FindByID(id uuid.UUID) (*domain.Item, error) {
	return itemService.repo.FindByID(id)
}
