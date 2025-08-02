package models

import (
	"github.com/google/uuid"
	"time"
)

type SwapRequestStatus string

const (
	StatusPending   SwapRequestStatus = "pending"
	StatusAccepted  SwapRequestStatus = "accepted"
	StatusRejected  SwapRequestStatus = "rejected"
	StatusCancelled SwapRequestStatus = "cancelled"
)

type SwapRequestModel struct {
	ID                   uuid.UUID         `gorm:"primary_key"`
	Status               SwapRequestStatus `gorm:"type:varchar(20)"`
	ReferenceNumber      string
	OfferedItemID        uuid.UUID `gorm:"type:uuid"`
	RequestedItemID      uuid.UUID `gorm:"type:uuid"`
	OfferedItemOwnerID   uuid.UUID `gorm:"type:uuid"`
	RequestedItemOwnerID uuid.UUID `gorm:"type:uuid"`
	CreatedAt            time.Time
	UpdatedAt            time.Time
}

func (SwapRequestModel) TableName() string {
	return "swapp_requests"
}
