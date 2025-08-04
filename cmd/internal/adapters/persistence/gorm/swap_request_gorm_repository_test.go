package gorm_test

import (
	"github.com/google/uuid"
	"github.com/stretchr/testify/assert"
	"gorm.io/gorm"
	gormRepo "swapp-go/cmd/internal/adapters/persistence/gorm"
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

func cleanSwapRequestsTable(t *testing.T, db *gorm.DB) {
	err := db.Session(&gorm.Session{AllowGlobalUpdate: true}).Delete(&models.SwapRequestModel{}).Error
	assert.NoError(t, err)
}

func TestSwapRequestGormRepository(t *testing.T) {
	db := testutils.SetupTestDB(t, &models.SwapRequestModel{})
	repo := gormRepo.NewSwapRequestGormRepository(db)

	t.Run("Create", func(t *testing.T) {
		swap := createTestSwapRequest(uuid.New(), uuid.New(), uuid.New(), uuid.New())
		err := repo.Create(swap)
		assert.NoError(t, err)
		assert.NotEqual(t, uuid.Nil, swap.ID)
	})

	t.Run("FindByID", func(t *testing.T) {
		swap := createTestSwapRequest(uuid.New(), uuid.New(), uuid.New(), uuid.New())
		err := repo.Create(swap)
		assert.NoError(t, err)

		fetched, err := repo.FindByID(swap.ID)
		assert.NoError(t, err)
		assert.Equal(t, swap.ID, fetched.ID)
	})

	t.Run("FindByReferenceNumber", func(t *testing.T) {
		cleanSwapRequestsTable(t, db)

		swap := createTestSwapRequest(uuid.New(), uuid.New(), uuid.New(), uuid.New())
		swap.ReferenceNumber = referenceNumber
		err := repo.Create(swap)
		assert.NoError(t, err)

		found, err := repo.FindByReferenceNumber(referenceNumber)
		assert.NoError(t, err)
		assert.Equal(t, swap.ID, found.ID)
	})

	t.Run("ListByUser", func(t *testing.T) {
		userA := uuid.New()
		userB := uuid.New()

		swap1 := createTestSwapRequest(uuid.New(), uuid.New(), userA, userB)
		swap2 := createTestSwapRequest(uuid.New(), uuid.New(), userB, userA)

		_ = repo.Create(swap1)
		_ = repo.Create(swap2)

		list, err := repo.ListByUser(userA)
		assert.NoError(t, err)
		assert.Len(t, list, 2)
	})

	t.Run("ListByStatus", func(t *testing.T) {
		cleanSwapRequestsTable(t, db)
		swap := createTestSwapRequest(uuid.New(), uuid.New(), uuid.New(), uuid.New())
		swap.Status = domain.StatusPending
		err := repo.Create(swap)
		assert.NoError(t, err)

		list, err := repo.ListByStatus(domain.StatusPending)
		assert.NoError(t, err)
		assert.Len(t, list, 1)
		assert.Equal(t, domain.StatusPending, list[0].Status)
	})

	t.Run("UpdateStatus", func(t *testing.T) {
		swap := createTestSwapRequest(uuid.New(), uuid.New(), uuid.New(), uuid.New())
		err := repo.Create(swap)
		assert.NoError(t, err)

		err = repo.UpdateStatus(swap.ID, domain.StatusAccepted)
		assert.NoError(t, err)

		updated, err := repo.FindByID(swap.ID)
		assert.NoError(t, err)
		assert.Equal(t, domain.StatusAccepted, updated.Status)
	})

	t.Run("Delete", func(t *testing.T) {
		swap := createTestSwapRequest(uuid.New(), uuid.New(), uuid.New(), uuid.New())
		err := repo.Create(swap)
		assert.NoError(t, err)

		err = repo.Delete(swap.ID)
		assert.NoError(t, err)

		_, err = repo.FindByID(swap.ID)
		assert.Error(t, err)
	})
}
