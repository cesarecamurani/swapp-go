package domain

import "github.com/google/uuid"

type SwappRequestStatus string

const (
	SwappStatusPending   SwappRequestStatus = "pending"
	SwappStatusAccepted  SwappRequestStatus = "accepted"
	SwappStatusRejected  SwappRequestStatus = "rejected"
	SwappStatusCancelled SwappRequestStatus = "cancelled"
)

type SwappRequest struct {
	ID                   uuid.UUID
	Status               SwappRequestStatus
	ReferenceNumber      string
	OfferedItemID        uuid.UUID
	RequestedItemID      uuid.UUID
	OfferedItemOwnerID   uuid.UUID
	RequestedItemOwnerID uuid.UUID
}
