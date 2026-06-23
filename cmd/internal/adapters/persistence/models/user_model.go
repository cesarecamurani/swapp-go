package models

import (
	"github.com/google/uuid"
	"time"
)

type UserModel struct {
	ID        uuid.UUID `gorm:"primaryKey"`
	Username  string    `gorm:"uniqueIndex;not null"`
	Password  string    `gorm:"not null"`
	Email     string    `gorm:"uniqueIndex;not null"`
	Phone     *string
	Address   *string
	CreatedAt time.Time
	UpdatedAt time.Time
}

func (UserModel) TableName() string {
	return "users"
}
