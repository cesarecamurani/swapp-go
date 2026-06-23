package services

import (
	"errors"
	"fmt"
	"github.com/google/uuid"
	"swapp-go/cmd/internal/application/ports"
	"swapp-go/cmd/internal/domain"
)

var ItemAlreadyOfferedErr = errors.New("item is already out for offer")

type SwapRequestService struct {
	repo         ports.SwapRequestRepository
	userRepo     ports.UserRepository
	itemRepo     ports.ItemRepository
	emailService ports.EmailService
}

func NewSwapRequestService(
	repo ports.SwapRequestRepository,
	userRepo ports.UserRepository,
	itemRepo ports.ItemRepository,
	emailService ports.EmailService,
) *SwapRequestService {
	return &SwapRequestService{
		repo:         repo,
		userRepo:     userRepo,
		itemRepo:     itemRepo,
		emailService: emailService,
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

	subject := fmt.Sprintf("New Swap Request Created (reference %v)", request.ReferenceNumber)

	service.sendEmailToUser(
		request.RecipientID,
		subject,
		fmt.Sprintf("You have a new swap request from %s", service.getUsernameSafe(request.SenderID)),
	)

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

	subject := fmt.Sprintf("Swap request with reference %s has been %s", swapRequest.ReferenceNumber, status)

	switch status {
	case domain.StatusAccepted:
		service.sendEmailToUser(
			swapRequest.SenderID,
			subject,
			fmt.Sprintf("Good news! Your swap request has been accepted by %s.", service.getUsernameSafe(swapRequest.RecipientID)),
		)
	case domain.StatusRejected:
		service.sendEmailToUser(
			swapRequest.SenderID,
			subject,
			fmt.Sprintf("Sorry, your swap request has been rejected by %s.", service.getUsernameSafe(swapRequest.RecipientID)),
		)
		if err = service.setItemOfferedStatus(swapRequest.OfferedItemID, false); err != nil {
			return fmt.Errorf("error releasing item after rejection: %w", err)
		}
	case domain.StatusCancelled:
		if err = service.setItemOfferedStatus(swapRequest.OfferedItemID, false); err != nil {
			return fmt.Errorf("error releasing item after cancellation: %w", err)
		}
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

	subject := fmt.Sprintf("Swap request with reference %s has been cancelled", swapRequest.ReferenceNumber)

	service.sendEmailToUser(
		swapRequest.RecipientID,
		subject,
		fmt.Sprintf("The swap request from %s has been cancelled.", service.getUsernameSafe(swapRequest.SenderID)),
	)

	return service.setItemOfferedStatus(swapRequest.OfferedItemID, false)
}

func (service *SwapRequestService) setItemOfferedStatus(itemID uuid.UUID, offered bool) error {
	_, err := service.itemRepo.Update(itemID, map[string]interface{}{
		"offered": offered,
	})

	return err
}

func (service *SwapRequestService) sendEmailToUser(userID uuid.UUID, subject, body string) {
	user, err := service.userRepo.FindByID(userID)
	if err != nil {
		fmt.Printf("Failed to find user %s for email: %v\n", userID, err)
		return
	}

	email := &domain.EmailMessage{
		Recipient: user.Email,
		Subject:   subject,
		Body:      body,
	}

	if err = service.emailService.SendEmail(email); err != nil {
		fmt.Printf("Failed to send email to %s: %v\n", user.Email, err)
	}
}

func (service *SwapRequestService) getUsernameSafe(userID uuid.UUID) string {
	user, err := service.userRepo.FindByID(userID)
	if err != nil {
		return "Unknown User"
	}
	return user.Username
}
