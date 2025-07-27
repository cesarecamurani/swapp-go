package service

import (
	"errors"
	"github.com/google/uuid"
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
	existingEmail, _ := userService.repo.GetUserByEmail(user.Email)
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

func (userService *UserService) UpdateUser(id uuid.UUID, fields map[string]interface{}) (*domain.User, error) {
	_, err := userService.repo.GetUserByID(id)
	if err != nil {
		return nil, err
	}

	updatedUser, err := userService.repo.UpdateUser(id, fields)
	if err != nil {
		return nil, err
	}

	return updatedUser, nil
}

func (userService *UserService) GetUserByID(id uuid.UUID) (*domain.User, error) {
	return userService.repo.GetUserByID(id)
}

func (userService *UserService) GetUserByUsername(username string) (*domain.User, error) {
	return userService.repo.GetUserByUsername(username)
}

func (userService *UserService) GetUserByEmail(email string) (*domain.User, error) {
	return userService.repo.GetUserByEmail(email)
}

func (userService *UserService) Authenticate(username, password string) (string, *domain.User, error) {
	user, err := userService.repo.GetUserByUsername(username)
	if err != nil {
		return "", nil, errors.New("invalid username")
	}

	if !utils.CheckPasswordHash(password, user.Password) {
		return "", nil, errors.New("invalid credentials")
	}

	token, err := utils.GenerateToken(user.Email, user.ID.String())
	if err != nil {
		return "", nil, err
	}

	return token, user, nil
}
