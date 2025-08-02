package services

import (
	"github.com/google/uuid"
	"swapp-go/cmd/internal/application/ports"
	"swapp-go/cmd/internal/domain"
)

type SwapRequestService struct {
	repo ports.SwapRequestRepository
}

func NewSwapRequestService(repo ports.SwapRequestRepository) *SwapRequestService {
	return &SwapRequestService{repo: repo}
}

func (service *SwapRequestService) Create(request *domain.SwapRequest) error {
	return service.repo.Create(request)
}

func (service *SwapRequestService) FindByID(id uuid.UUID) (*domain.SwapRequest, error) {
	return service.repo.FindByID(id)
}

func (service *SwapRequestService) FindByReferenceNumber(reference string) (*domain.SwapRequest, error) {
	return service.repo.FindByReferenceNumber(reference)
}

func (service *SwapRequestService) ListByUser(userID uuid.UUID) ([]domain.SwapRequest, error) {
	return service.repo.ListByUser(userID)
}

func (service *SwapRequestService) ListByStatus(status domain.SwapRequestStatus) ([]domain.SwapRequest, error) {
	return service.repo.ListByStatus(status)
}

func (service *SwapRequestService) UpdateStatus(id uuid.UUID, status domain.SwapRequestStatus) error {
	return service.repo.UpdateStatus(id, status)
}

func (service *SwapRequestService) Delete(id uuid.UUID) error {
	return service.repo.Delete(id)
}
