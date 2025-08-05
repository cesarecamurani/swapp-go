package mocks

import (
	"github.com/google/uuid"
	"github.com/stretchr/testify/mock"
	"swapp-go/cmd/internal/domain"
)

type SwapRequestRepository struct {
	mock.Mock
}

func (m *SwapRequestRepository) Create(request *domain.SwapRequest) error {
	args := m.Called(request)
	return args.Error(0)
}

func (m *SwapRequestRepository) FindByID(id uuid.UUID) (*domain.SwapRequest, error) {
	args := m.Called(id)
	if req, ok := args.Get(0).(*domain.SwapRequest); ok {
		return req, args.Error(1)
	}
	return nil, args.Error(1)
}

func (m *SwapRequestRepository) FindByReferenceNumber(reference string) (*domain.SwapRequest, error) {
	args := m.Called(reference)
	if req, ok := args.Get(0).(*domain.SwapRequest); ok {
		return req, args.Error(1)
	}
	return nil, args.Error(1)
}

func (m *SwapRequestRepository) ListByUser(userID uuid.UUID) ([]domain.SwapRequest, error) {
	args := m.Called(userID)
	if list, ok := args.Get(0).([]domain.SwapRequest); ok {
		return list, args.Error(1)
	}
	return nil, args.Error(1)
}

func (m *SwapRequestRepository) ListByStatus(status domain.SwapRequestStatus) ([]domain.SwapRequest, error) {
	args := m.Called(status)
	if list, ok := args.Get(0).([]domain.SwapRequest); ok {
		return list, args.Error(1)
	}
	return nil, args.Error(1)
}

func (m *SwapRequestRepository) UpdateStatus(id uuid.UUID, status domain.SwapRequestStatus) error {
	args := m.Called(id, status)
	return args.Error(0)
}

func (m *SwapRequestRepository) Delete(id uuid.UUID) error {
	args := m.Called(id)
	return args.Error(0)
}
