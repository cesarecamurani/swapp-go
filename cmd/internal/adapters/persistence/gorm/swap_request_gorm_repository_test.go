package gorm_test

import (
	"github.com/google/uuid"
	"github.com/stretchr/testify/assert"
	"swapp-go/cmd/internal/adapters/persistence/gorm"
	"swapp-go/cmd/internal/adapters/persistence/gorm/testutils"
	"swapp-go/cmd/internal/adapters/persistence/models"
	"swapp-go/cmd/internal/domain"
	"testing"
)

var referenceNumber = "REF123456"

func createTestSwapRequest(offeredItemID, requestedItemID, senderID, receiverID uuid.UUID) *domain.SwapRequest {
	return &domain.SwapRequest{
		Status:          domain.StatusPending,
		ReferenceNumber: referenceNumber,
		OfferedItemID:   offeredItemID,
		RequestedItemID: requestedItemID,
		SenderID:        senderID,
		RecipientID:     receiverID,
	}
}

func TestSwapRequestGormRepository_Create(t *testing.T) {
	db := testutils.SetupTestDB(t, &models.SwapRequestModel{})
	repository := gorm.NewSwapRequestGormRepository(db)

	swapRequest := createTestSwapRequest(uuid.New(), uuid.New(), uuid.New(), uuid.New())
	err := repository.Create(swapRequest)
	assert.NoError(t, err)
	assert.NotEqual(t, uuid.Nil, swapRequest.ID)
}

func TestSwapRequestGormRepository_FindByID(t *testing.T) {
	db := testutils.SetupTestDB(t, &models.SwapRequestModel{})
	repository := gorm.NewSwapRequestGormRepository(db)

	swapRequest := createTestSwapRequest(uuid.New(), uuid.New(), uuid.New(), uuid.New())
	err := repository.Create(swapRequest)
	assert.NoError(t, err)

	fetched, err := repository.FindByID(swapRequest.ID)
	assert.NoError(t, err)
	assert.Equal(t, swapRequest.ID, fetched.ID)
}

func TestSwapRequestGormRepository_FindByReferenceNumber(t *testing.T) {
	db := testutils.SetupTestDB(t, &models.SwapRequestModel{})
	repository := gorm.NewSwapRequestGormRepository(db)

	swapRequest := createTestSwapRequest(uuid.New(), uuid.New(), uuid.New(), uuid.New())
	swapRequest.ReferenceNumber = referenceNumber
	_ = repository.Create(swapRequest)

	found, err := repository.FindByReferenceNumber(referenceNumber)
	assert.NoError(t, err)
	assert.Equal(t, swapRequest.ID, found.ID)
}

func TestSwapRequestGormRepository_ListByUser(t *testing.T) {
	db := testutils.SetupTestDB(t, &models.SwapRequestModel{})
	repository := gorm.NewSwapRequestGormRepository(db)

	userA := uuid.New()
	userB := uuid.New()

	swap1 := createTestSwapRequest(uuid.New(), uuid.New(), userA, userB)
	swap2 := createTestSwapRequest(uuid.New(), uuid.New(), userB, userA)
	_ = repository.Create(swap1)
	_ = repository.Create(swap2)

	list, err := repository.ListByUser(userA)
	assert.NoError(t, err)
	assert.Len(t, list, 2)
}

func TestSwapRequestGormRepository_ListByStatus(t *testing.T) {
	db := testutils.SetupTestDB(t, &models.SwapRequestModel{})
	repository := gorm.NewSwapRequestGormRepository(db)

	swap := createTestSwapRequest(uuid.New(), uuid.New(), uuid.New(), uuid.New())
	swap.Status = domain.StatusPending
	_ = repository.Create(swap)

	list, err := repository.ListByStatus(domain.StatusPending)
	assert.NoError(t, err)
	assert.Len(t, list, 1)
	assert.Equal(t, domain.StatusPending, list[0].Status)
}

func TestSwapRequestGormRepository_UpdateStatus(t *testing.T) {
	db := testutils.SetupTestDB(t, &models.SwapRequestModel{})
	repository := gorm.NewSwapRequestGormRepository(db)

	swap := createTestSwapRequest(uuid.New(), uuid.New(), uuid.New(), uuid.New())
	_ = repository.Create(swap)

	err := repository.UpdateStatus(swap.ID, domain.StatusAccepted)
	assert.NoError(t, err)

	updated, err := repository.FindByID(swap.ID)
	assert.NoError(t, err)
	assert.Equal(t, domain.StatusAccepted, updated.Status)
}

func TestSwapRequestGormRepository_Delete(t *testing.T) {
	db := testutils.SetupTestDB(t, &models.SwapRequestModel{})
	repository := gorm.NewSwapRequestGormRepository(db)

	swap := createTestSwapRequest(uuid.New(), uuid.New(), uuid.New(), uuid.New())
	_ = repository.Create(swap)

	err := repository.Delete(swap.ID)
	assert.NoError(t, err)

	_, err = repository.FindByID(swap.ID)
	assert.Error(t, err)
}
