package services

import (
	"github.com/google/uuid"
	"swapp-go/cmd/internal/domain"
)

type ItemServiceInterface interface {
	CreateItem(item *domain.Item) error
	UpdateItem(id uuid.UUID, fields map[string]interface{}) (*domain.Item, error)
	DeleteItem(id uuid.UUID) error
	GetItemByID(id uuid.UUID) (*domain.Item, error)
}
