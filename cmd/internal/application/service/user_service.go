package service

import (
	"errors"
	"swapp-go/cmd/internal/application/ports"
	"swapp-go/cmd/internal/domain"
	"swapp-go/cmd/internal/utils"
)

type UserService struct {
	repo ports.UserRepository
}

func NewUserService(repo ports.UserRepository) *UserService {
	return &UserService{repo: repo}
}

func (userService *UserService) RegisterUser(user *domain.User) error {
	existingEmail, _ := userService.repo.GetUserByEmail(user.Username)
	if existingEmail != nil {
		return errors.New("email already exists")
	}

	existingUsername, _ := userService.repo.GetUserByUsername(user.Username)
	if existingUsername != nil {
		return errors.New("username not available")
	}

	encryptedPassword, err := utils.HashPassword(user.Password)
	if err != nil {
		return err
	}

	user.Password = encryptedPassword

	return userService.repo.CreateUser(user)
}

func (userService *UserService) GetUserByID(id uint) (*domain.User, error) {
	return userService.repo.GetUserByID(id)
}

func (userService *UserService) GetUserByUsername(username string) (*domain.User, error) {
	return userService.repo.GetUserByUsername(username)
}
