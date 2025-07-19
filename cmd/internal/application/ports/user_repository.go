package ports

import "swapp-go/cmd/internal/domain"

type UserRepository interface {
	CreateUser(user *domain.User) error
	GetUserByID(id uint) (*domain.User, error)
	GetUserByUsername(username string) (*domain.User, error)
	GetUserByEmail(email string) (*domain.User, error)
}
