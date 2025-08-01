package gorm

import (
	"github.com/google/uuid"
	"gorm.io/gorm"
	"swapp-go/cmd/internal/adapters/persistence/models"
	"swapp-go/cmd/internal/domain"
)

type SwappRequestGormRepository struct {
	db *gorm.DB
}

func NewSwappRequestGormRepository(db *gorm.DB) *SwappRequestGormRepository {
	return &SwappRequestGormRepository{db}
}

func toSwappRequestModel(swappRequest *domain.SwappRequest) *models.SwappRequestModel {
	id := swappRequest.ID
	if id == uuid.Nil {
		id = uuid.New()
	}

	return &models.SwappRequestModel{
		ID:                   id,
		Status:               models.SwappRequestStatus(swappRequest.Status),
		ReferenceNumber:      swappRequest.ReferenceNumber,
		OfferedItemID:        swappRequest.OfferedItemID,
		RequestedItemID:      swappRequest.RequestedItemID,
		OfferedItemOwnerID:   swappRequest.OfferedItemOwnerID,
		RequestedItemOwnerID: swappRequest.RequestedItemOwnerID,
	}
}

func toDomainSwappRequest(model *models.SwappRequestModel) *domain.SwappRequest {
	return &domain.SwappRequest{
		ID:                   model.ID,
		Status:               domain.SwappRequestStatus(model.Status),
		ReferenceNumber:      model.ReferenceNumber,
		OfferedItemID:        model.OfferedItemID,
		RequestedItemID:      model.RequestedItemID,
		OfferedItemOwnerID:   model.OfferedItemOwnerID,
		RequestedItemOwnerID: model.RequestedItemOwnerID,
	}
}

func (swappRequestGorm *SwappRequestGormRepository) Create(swappRequest *domain.SwappRequest) error {
	model := toSwappRequestModel(swappRequest)

	if result := swappRequestGorm.db.Create(model); result.Error != nil {
		return result.Error
	}

	swappRequest.ID = model.ID

	return nil
}

func (swappRequestGorm *SwappRequestGormRepository) FindByID(id uuid.UUID) (*domain.SwappRequest, error) {
	var model models.SwappRequestModel
	if err := swappRequestGorm.db.First(&model, "id = ?", id).Error; err != nil {
		return nil, err
	}
	return toDomainSwappRequest(&model), nil
}

func (swappRequestGorm *SwappRequestGormRepository) FindByReferenceNumber(reference string) (*domain.SwappRequest, error) {
	var model models.SwappRequestModel
	if err := swappRequestGorm.db.First(&model, "reference_number = ?", reference).Error; err != nil {
		return nil, err
	}
	return toDomainSwappRequest(&model), nil
}

func (swappRequestGorm *SwappRequestGormRepository) ListByUser(userID uuid.UUID) ([]domain.SwappRequest, error) {
	var modelsList []models.SwappRequestModel
	if err := swappRequestGorm.db.Where(
		"offered_item_owner_id = ? OR requested_item_owner_id = ?", userID, userID,
	).Find(&modelsList).Error; err != nil {
		return nil, err
	}

	var domainList []domain.SwappRequest
	for _, m := range modelsList {
		domainList = append(domainList, *toDomainSwappRequest(&m))
	}

	return domainList, nil
}

func (swappRequestGorm *SwappRequestGormRepository) ListByStatus(status domain.SwappRequestStatus) ([]domain.SwappRequest, error) {
	var modelsList []models.SwappRequestModel
	if err := swappRequestGorm.db.Where("status = ?", string(status)).Find(&modelsList).Error; err != nil {
		return nil, err
	}

	var domainList []domain.SwappRequest
	for _, m := range modelsList {
		domainList = append(domainList, *toDomainSwappRequest(&m))
	}

	return domainList, nil
}

func (swappRequestGorm *SwappRequestGormRepository) UpdateStatus(id uuid.UUID, status domain.SwappRequestStatus) error {
	return swappRequestGorm.db.Model(&models.SwappRequestModel{}).
		Where("id = ?", id).
		Update("status", string(status)).Error
}

func (swappRequestGorm *SwappRequestGormRepository) Delete(id uuid.UUID) error {
	return swappRequestGorm.db.Delete(&models.SwappRequestModel{}, "id = ?", id).Error
}
