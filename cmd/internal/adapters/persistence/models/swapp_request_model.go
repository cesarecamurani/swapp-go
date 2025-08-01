package models

import (
	"github.com/google/uuid"
	"time"
)

type SwappRequestStatus string

const (
	SwappStatusPending   SwappRequestStatus = "pending"
	SwappStatusAccepted  SwappRequestStatus = "accepted"
	SwappStatusRejected  SwappRequestStatus = "rejected"
	SwappStatusCancelled SwappRequestStatus = "cancelled"
)

type SwappRequestModel struct {
	ID                   uuid.UUID          `gorm:"primary_key"`
	Status               SwappRequestStatus `gorm:"type:varchar(20)"`
	ReferenceNumber      string
	OfferedItemID        uuid.UUID `gorm:"type:uuid"`
	RequestedItemID      uuid.UUID `gorm:"type:uuid"`
	OfferedItemOwnerID   uuid.UUID `gorm:"type:uuid"`
	RequestedItemOwnerID uuid.UUID `gorm:"type:uuid"`
	CreatedAt            time.Time
	UpdatedAt            time.Time
}

func (SwappRequestModel) TableName() string {
	return "swapp_requests"
}
