package gorm

import (
	"github.com/google/uuid"
	"gorm.io/gorm"
	"swapp-go/cmd/internal/adapters/persistence/models"
	"swapp-go/cmd/internal/domain"
)

type SwapRequestGormRepository struct {
	db *gorm.DB
}

func NewSwapRequestGormRepository(db *gorm.DB) *SwapRequestGormRepository {
	return &SwapRequestGormRepository{db}
}

func toSwapRequestModel(swapRequest *domain.SwapRequest) *models.SwapRequestModel {
	id := swapRequest.ID
	if id == uuid.Nil {
		id = uuid.New()
	}

	return &models.SwapRequestModel{
		ID:                   id,
		Status:               models.SwapRequestStatus(swapRequest.Status),
		ReferenceNumber:      swapRequest.ReferenceNumber,
		OfferedItemID:        swapRequest.OfferedItemID,
		RequestedItemID:      swapRequest.RequestedItemID,
		OfferedItemOwnerID:   swapRequest.OfferedItemOwnerID,
		RequestedItemOwnerID: swapRequest.RequestedItemOwnerID,
	}
}

func toDomainSwapRequest(model *models.SwapRequestModel) *domain.SwapRequest {
	return &domain.SwapRequest{
		ID:                   model.ID,
		Status:               domain.SwapRequestStatus(model.Status),
		ReferenceNumber:      model.ReferenceNumber,
		OfferedItemID:        model.OfferedItemID,
		RequestedItemID:      model.RequestedItemID,
		OfferedItemOwnerID:   model.OfferedItemOwnerID,
		RequestedItemOwnerID: model.RequestedItemOwnerID,
	}
}

func (swapRequestGorm *SwapRequestGormRepository) Create(swapRequest *domain.SwapRequest) error {
	model := toSwapRequestModel(swapRequest)

	if result := swapRequestGorm.db.Create(model); result.Error != nil {
		return result.Error
	}

	swapRequest.ID = model.ID

	return nil
}

func (swapRequestGorm *SwapRequestGormRepository) FindByID(id uuid.UUID) (*domain.SwapRequest, error) {
	var model models.SwapRequestModel
	if err := swapRequestGorm.db.First(&model, "id = ?", id).Error; err != nil {
		return nil, err
	}
	return toDomainSwapRequest(&model), nil
}

func (swapRequestGorm *SwapRequestGormRepository) FindByReferenceNumber(reference string) (*domain.SwapRequest, error) {
	var model models.SwapRequestModel
	if err := swapRequestGorm.db.First(&model, "reference_number = ?", reference).Error; err != nil {
		return nil, err
	}
	return toDomainSwapRequest(&model), nil
}

func (swapRequestGorm *SwapRequestGormRepository) ListByUser(userID uuid.UUID) ([]domain.SwapRequest, error) {
	var modelsList []models.SwapRequestModel
	if err := swapRequestGorm.db.Where(
		"offered_item_owner_id = ? OR requested_item_owner_id = ?", userID, userID,
	).Find(&modelsList).Error; err != nil {
		return nil, err
	}

	var domainList []domain.SwapRequest
	for _, m := range modelsList {
		domainList = append(domainList, *toDomainSwapRequest(&m))
	}

	return domainList, nil
}

func (swapRequestGorm *SwapRequestGormRepository) ListByStatus(status domain.SwapRequestStatus) ([]domain.SwapRequest, error) {
	var modelsList []models.SwapRequestModel
	if err := swapRequestGorm.db.Where("status = ?", string(status)).Find(&modelsList).Error; err != nil {
		return nil, err
	}

	var domainList []domain.SwapRequest
	for _, m := range modelsList {
		domainList = append(domainList, *toDomainSwapRequest(&m))
	}

	return domainList, nil
}

func (swapRequestGorm *SwapRequestGormRepository) UpdateStatus(id uuid.UUID, status domain.SwapRequestStatus) error {
	return swapRequestGorm.db.Model(&models.SwapRequestModel{}).
		Where("id = ?", id).
		Update("status", string(status)).Error
}

func (swapRequestGorm *SwapRequestGormRepository) Delete(id uuid.UUID) error {
	return swapRequestGorm.db.Delete(&models.SwapRequestModel{}, "id = ?", id).Error
}
