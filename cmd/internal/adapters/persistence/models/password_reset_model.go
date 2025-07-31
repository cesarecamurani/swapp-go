package models

import (
	"github.com/google/uuid"
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
