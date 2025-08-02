package domain

import "github.com/google/uuid"

type SwapRequestStatus string

const (
	StatusPending   SwapRequestStatus = "pending"
	StatusAccepted  SwapRequestStatus = "accepted"
	StatusRejected  SwapRequestStatus = "rejected"
	StatusCancelled SwapRequestStatus = "cancelled"
)

type SwapRequest struct {
	ID                   uuid.UUID
	Status               SwapRequestStatus
	ReferenceNumber      string
	OfferedItemID        uuid.UUID
	RequestedItemID      uuid.UUID
	OfferedItemOwnerID   uuid.UUID
	RequestedItemOwnerID uuid.UUID
}
