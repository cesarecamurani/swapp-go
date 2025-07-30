package persistence

import (
	"github.com/google/uuid"
	"gorm.io/gorm"
	"swapp-go/cmd/internal/domain"
	"time"
)

type PasswordResetModel struct {
	Token     string    `gorm:"primaryKey"`
	UserID    uuid.UUID `gorm:"type:uuid;not null"`
	ExpiresAt time.Time `gorm:"not null"`
}

func (PasswordResetModel) TableName() string {
	return "password_resets"
}

type GormPasswordResetRepository struct {
	db *gorm.DB
}

func NewGormPasswordResetRepository(db *gorm.DB) *GormPasswordResetRepository {
	return &GormPasswordResetRepository{db: db}
}

func (r *GormPasswordResetRepository) Save(token *domain.PasswordReset) error {
	model := PasswordResetModel{
		Token:     token.Token,
		UserID:    token.UserID,
		ExpiresAt: token.ExpiresAt,
	}

	return r.db.Create(&model).Error
}

func (r *GormPasswordResetRepository) GetByToken(token string) (*domain.PasswordReset, error) {
	var model PasswordResetModel

	if err := r.db.First(&model, "token = ?", token).Error; err != nil {
		return nil, err
	}

	return &domain.PasswordReset{
		Token:     model.Token,
		UserID:    model.UserID,
		ExpiresAt: model.ExpiresAt,
	}, nil
}

func (r *GormPasswordResetRepository) Delete(token string) error {
	return r.db.Delete(&PasswordResetModel{}, "token = ?", token).Error
}
