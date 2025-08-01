package services

import (
	"github.com/google/uuid"
	"swapp-go/cmd/internal/application/ports"
	"swapp-go/cmd/internal/domain"
)

type SwappRequestService struct {
	repo ports.SwappRequestRepository
}

func NewSwappRequestService(repo ports.SwappRequestRepository) *SwappRequestService {
	return &SwappRequestService{repo: repo}
}

func (service *SwappRequestService) CreateSwappRequest(request *domain.SwappRequest) error {
	return service.repo.Create(request)
}

func (service *SwappRequestService) GetSwappRequestByID(id uuid.UUID) (*domain.SwappRequest, error) {
	return service.repo.FindByID(id)
}

func (service *SwappRequestService) GetSwappRequestByReferenceNumber(reference string) (*domain.SwappRequest, error) {
	return service.repo.FindByReferenceNumber(reference)
}

func (service *SwappRequestService) ListSwappRequestsByUser(userID uuid.UUID) ([]domain.SwappRequest, error) {
	return service.repo.ListByUser(userID)
}

func (service *SwappRequestService) ListSwappRequestsByStatus(status domain.SwappRequestStatus) ([]domain.SwappRequest, error) {
	return service.repo.ListByStatus(status)
}

func (service *SwappRequestService) UpdateSwappRequestStatus(id uuid.UUID, status domain.SwappRequestStatus) error {
	return service.repo.UpdateStatus(id, status)
}

func (service *SwappRequestService) DeleteSwappRequest(id uuid.UUID) error {
	return service.repo.Delete(id)
}
