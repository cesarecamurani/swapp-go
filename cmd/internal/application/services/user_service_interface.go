package services

import (
	"github.com/google/uuid"
	"swapp-go/cmd/internal/domain"
)

type UserServiceInterface interface {
	RegisterUser(user *domain.User) error
	Update(id uuid.UUID, fields map[string]interface{}) (*domain.User, error)
	Delete(id uuid.UUID) error
	FindByID(id uuid.UUID) (*domain.User, error)
	FindByEmail(email string) (*domain.User, error)
	FindByUsername(username string) (*domain.User, error)
	Authenticate(username, password string) (string, *domain.User, error)
}
