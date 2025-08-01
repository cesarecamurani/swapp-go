package ports

import (
	"github.com/google/uuid"
	"swapp-go/cmd/internal/domain"
)

type SwappRequestRepository interface {
	Create(request *domain.SwappRequest) error
	FindByID(id uuid.UUID) (*domain.SwappRequest, error)
	FindByReferenceNumber(reference string) (*domain.SwappRequest, error)
	ListByUser(userID uuid.UUID) ([]domain.SwappRequest, error)
	ListByStatus(status domain.SwappRequestStatus) ([]domain.SwappRequest, error)
	UpdateStatus(id uuid.UUID, status domain.SwappRequestStatus) error
	Delete(id uuid.UUID) error
}
