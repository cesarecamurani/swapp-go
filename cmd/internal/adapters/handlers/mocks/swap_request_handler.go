package mocks

import (
	"github.com/google/uuid"
	"github.com/stretchr/testify/mock"
	"swapp-go/cmd/internal/domain"
)

type SwapRequestService struct {
	mock.Mock
}

func (m *SwapRequestService) Create(request *domain.SwapRequest) error {
	args := m.Called(request)
	return args.Error(0)
}

func (m *SwapRequestService) FindByID(id uuid.UUID) (*domain.SwapRequest, error) {
	args := m.Called(id)
	result := args.Get(0)
	if result == nil {
		return nil, args.Error(1)
	}
	return result.(*domain.SwapRequest), args.Error(1)
}

func (m *SwapRequestService) FindByReferenceNumber(reference string) (*domain.SwapRequest, error) {
	args := m.Called(reference)
	result := args.Get(0)
	if result == nil {
		return nil, args.Error(1)
	}
	return result.(*domain.SwapRequest), args.Error(1)
}

func (m *SwapRequestService) ListByUser(userID uuid.UUID) ([]domain.SwapRequest, error) {
	args := m.Called(userID)
	result := args.Get(0)
	if result == nil {
		return nil, args.Error(1)
	}
	return result.([]domain.SwapRequest), args.Error(1)
}

func (m *SwapRequestService) ListByStatus(status domain.SwapRequestStatus) ([]domain.SwapRequest, error) {
	args := m.Called(status)
	result := args.Get(0)
	if result == nil {
		return nil, args.Error(1)
	}
	return result.([]domain.SwapRequest), args.Error(1)
}

func (m *SwapRequestService) UpdateStatus(id uuid.UUID, status domain.SwapRequestStatus) error {
	args := m.Called(id, status)
	return args.Error(0)
}

func (m *SwapRequestService) Delete(id uuid.UUID) error {
	args := m.Called(id)
	return args.Error(0)
}
