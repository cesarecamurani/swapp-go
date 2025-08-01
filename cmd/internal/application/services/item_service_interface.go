package services

import (
	"github.com/google/uuid"
	"swapp-go/cmd/internal/domain"
)

type ItemServiceInterface interface {
	Create(item *domain.Item) error
	Update(id uuid.UUID, fields map[string]interface{}) (*domain.Item, error)
	Delete(id uuid.UUID) error
	FindByID(id uuid.UUID) (*domain.Item, error)
}
