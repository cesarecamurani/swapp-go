package gorm

import (
	"gorm.io/gorm"
	"swapp-go/cmd/internal/adapters/persistence/models"
	"swapp-go/cmd/internal/domain"
)

type PasswordResetGormRepository struct {
	db *gorm.DB
}

func NewPasswordResetGormRepository(db *gorm.DB) *PasswordResetGormRepository {
	return &PasswordResetGormRepository{db: db}
}

func (r *PasswordResetGormRepository) Save(token *domain.PasswordReset) error {
	model := models.PasswordResetModel{
		Token:     token.Token,
		UserID:    token.UserID,
		ExpiresAt: token.ExpiresAt,
	}

	return r.db.Create(&model).Error
}

func (r *PasswordResetGormRepository) GetByToken(token string) (*domain.PasswordReset, error) {
	var model models.PasswordResetModel

	if err := r.db.First(&model, "token = ?", token).Error; err != nil {
		return nil, err
	}

	return &domain.PasswordReset{
		Token:     model.Token,
		UserID:    model.UserID,
		ExpiresAt: model.ExpiresAt,
	}, nil
}

func (r *PasswordResetGormRepository) Delete(token string) error {
	return r.db.Delete(&models.PasswordResetModel{}, "token = ?", token).Error
}
