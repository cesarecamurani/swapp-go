package ports

import (
	"github.com/google/uuid"
	"swapp-go/cmd/internal/domain"
)

type SwapRequestRepository interface {
	Create(request *domain.SwapRequest) error
	FindByID(id uuid.UUID) (*domain.SwapRequest, error)
	FindByReferenceNumber(reference string) (*domain.SwapRequest, error)
	ListByUser(userID uuid.UUID) ([]domain.SwapRequest, error)
	ListByStatus(status domain.SwapRequestStatus) ([]domain.SwapRequest, error)
	UpdateStatus(id uuid.UUID, status domain.SwapRequestStatus) error
	Delete(id uuid.UUID) error
}
