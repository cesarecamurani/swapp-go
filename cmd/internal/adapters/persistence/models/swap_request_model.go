package models

import (
	"github.com/google/uuid"
	"time"
)

type SwapRequestStatus string

type SwapRequestModel struct {
	ID              uuid.UUID         `gorm:"primary_key"`
	Status          SwapRequestStatus `gorm:"type:varchar(20)"`
	ReferenceNumber string
	OfferedItemID   uuid.UUID `gorm:"type:uuid"`
	RequestedItemID uuid.UUID `gorm:"type:uuid"`
	SenderID        uuid.UUID `gorm:"type:uuid"`
	RecipientID     uuid.UUID `gorm:"type:uuid"`
	CreatedAt       time.Time
	UpdatedAt       time.Time
}

func (SwapRequestModel) TableName() string {
	return "swap_requests"
}
