package ports

import "swapp-go/cmd/internal/domain"

type PasswordResetRepository interface {
	Save(token *domain.PasswordReset) error
	GetByToken(token string) (*domain.PasswordReset, error)
	Delete(token string) error
}
