package ports

import (
	"github.com/google/uuid"
	"swapp-go/cmd/internal/domain"
)

type UserRepository interface {
	CreateUser(user *domain.User) error
	UpdateUser(id uuid.UUID, fields map[string]interface{}) (*domain.User, error)
	GetUserByID(id uuid.UUID) (*domain.User, error)
	GetUserByUsername(username string) (*domain.User, error)
	GetUserByEmail(email string) (*domain.User, error)
}
