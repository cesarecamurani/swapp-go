package services

import (
	"github.com/google/uuid"
	"swapp-go/cmd/internal/domain"
)

type SwappRequestServiceInterface interface {
	CreateSwappRequest(request domain.SwappRequest) error
	GetSwappRequestByID(id uuid.UUID) (*domain.SwappRequest, error)
	DeleteSwappRequest(id uuid.UUID) error
	ListSwappRequestsByUser(userID uuid.UUID) ([]domain.SwappRequest, error)
	ListSwappRequestsByStatus(status domain.SwappRequestStatus) ([]domain.SwappRequest, error)
	UpdateSwappRequestStatus(id uuid.UUID, status domain.SwappRequestStatus) error
	GetSwappRequestByReferenceNumber(reference string) (*domain.SwappRequest, error)
}
