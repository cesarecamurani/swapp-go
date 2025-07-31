package services

import (
	"github.com/google/uuid"
	"swapp-go/cmd/internal/domain"
)

type UserServiceInterface interface {
	RegisterUser(user *domain.User) error
	UpdateUser(id uuid.UUID, fields map[string]interface{}) (*domain.User, error)
	DeleteUser(id uuid.UUID) error
	GetUserByID(id uuid.UUID) (*domain.User, error)
	GetUserByEmail(email string) (*domain.User, error)
	GetUserByUsername(username string) (*domain.User, error)
	Authenticate(username, password string) (string, *domain.User, error)
}
