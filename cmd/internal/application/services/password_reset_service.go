package services

import (
	"github.com/google/uuid"
	"swapp-go/cmd/internal/application/ports"
	"swapp-go/cmd/internal/domain"
	"time"
)

type PasswordResetService struct {
	ResetTokenRepo ports.PasswordResetRepository
}

func NewPasswordResetService(tokenRepo ports.PasswordResetRepository) *PasswordResetService {
	return &PasswordResetService{
		ResetTokenRepo: tokenRepo,
	}
}

func (service *PasswordResetService) GenerateAndSaveToken(userID uuid.UUID) (string, error) {
	token := uuid.NewString()
	resetToken := &domain.PasswordReset{
		Token:     token,
		UserID:    userID,
		ExpiresAt: time.Now().Add(1 * time.Hour),
	}

	err := service.ResetTokenRepo.Save(resetToken)
	if err != nil {
		return "", err
	}

	return token, nil
}

func (service *PasswordResetService) ValidateToken(token string) (*domain.PasswordReset, error) {
	return service.ResetTokenRepo.GetByToken(token)
}

func (service *PasswordResetService) DeleteToken(token string) error {
	return service.ResetTokenRepo.Delete(token)
}
