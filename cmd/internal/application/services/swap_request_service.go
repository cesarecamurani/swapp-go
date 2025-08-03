package services

import (
	"errors"
	"github.com/google/uuid"
	"swapp-go/cmd/internal/application/ports"
	"swapp-go/cmd/internal/domain"
)

var ItemAlreadyOfferedErr = errors.New("item is already out for offer")

type SwapRequestService struct {
	repo     ports.SwapRequestRepository
	itemRepo ports.ItemRepository
}

func NewSwapRequestService(repo ports.SwapRequestRepository, itemRepo ports.ItemRepository) *SwapRequestService {
	return &SwapRequestService{
		repo:     repo,
		itemRepo: itemRepo,
	}
}

func (service *SwapRequestService) Create(request *domain.SwapRequest) error {
	offeredItemID := request.OfferedItemID

	item, err := service.itemRepo.FindByID(offeredItemID)
	if err != nil {
		return errors.New("offered item not found")
	}

	success, err := service.itemRepo.TryMarkItemAsOffered(item.ID)
	if err != nil {
		return err
	}
	if !success {
		return ItemAlreadyOfferedErr
	}

	if err = service.setItemOfferedStatus(item.ID, true); err != nil {
		return errors.New("failed to mark item as offered")
	}

	if err = service.repo.Create(request); err != nil {
		_ = service.setItemOfferedStatus(item.ID, false)
		return err
	}

	return nil
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
	swapRequest, err := service.repo.FindByID(id)
	if err != nil {
		return err
	}

	if err = service.repo.UpdateStatus(id, status); err != nil {
		return err
	}

	if status == domain.StatusRejected || status == domain.StatusCancelled {
		return service.setItemOfferedStatus(swapRequest.OfferedItemID, false)
	}
	
	return nil
}

func (service *SwapRequestService) Delete(id uuid.UUID) error {
	swapRequest, err := service.repo.FindByID(id)
	if err != nil {
		return err
	}

	if err = service.repo.Delete(id); err != nil {
		return err
	}

	return service.setItemOfferedStatus(swapRequest.OfferedItemID, false)
}

func (service *SwapRequestService) setItemOfferedStatus(itemID uuid.UUID, offered bool) error {
	_, err := service.itemRepo.Update(itemID, map[string]interface{}{
		"offered": offered,
	})

	return err
}
