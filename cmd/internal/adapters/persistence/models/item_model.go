package models

import (
	"github.com/google/uuid"
	"time"
)

type ItemModel struct {
	ID          uuid.UUID `gorm:"primary_key"`
	Name        string    `gorm:"not null"`
	Description string    `gorm:"not null"`
	PictureURL  string    `gorm:"not null"`
	UserID      uuid.UUID
	Offered     bool `gorm:"default:false"`
	CreatedAt   time.Time
	UpdatedAt   time.Time
}

func (ItemModel) TableName() string {
	return "items"
}
