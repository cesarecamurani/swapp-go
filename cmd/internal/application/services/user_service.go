package services

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
	existingEmail, _ := userService.repo.FindByEmail(user.Email)
	if existingEmail != nil {
		return errors.New("email already exists")
	}

	existingUsername, _ := userService.repo.FindByUsername(user.Username)
	if existingUsername != nil {
		return errors.New("username not available")
	}

	encryptedPassword, err := utils.HashPassword(user.Password)
	if err != nil {
		return err
	}

	user.Password = encryptedPassword

	return userService.repo.Create(user)
}

func (userService *UserService) Update(id uuid.UUID, fields map[string]interface{}) (*domain.User, error) {
	_, err := userService.repo.FindByID(id)
	if err != nil {
		return nil, err
	}

	updatedUser, err := userService.repo.Update(id, fields)
	if err != nil {
		return nil, err
	}

	return updatedUser, nil
}

func (userService *UserService) Delete(id uuid.UUID) error {
	return userService.repo.Delete(id)
}

func (userService *UserService) FindByID(id uuid.UUID) (*domain.User, error) {
	return userService.repo.FindByID(id)
}

func (userService *UserService) FindByUsername(username string) (*domain.User, error) {
	return userService.repo.FindByUsername(username)
}

func (userService *UserService) FindByEmail(email string) (*domain.User, error) {
	return userService.repo.FindByEmail(email)
}

func (userService *UserService) Authenticate(username, password string) (string, *domain.User, error) {
	user, err := userService.repo.FindByUsername(username)
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
