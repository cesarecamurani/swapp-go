package domain

import (
	"github.com/google/uuid"
	"time"
)

type Item struct {
	ID          uuid.UUID
	Name        string
	Description string
	PictureURL  string
	UserID      uuid.UUID
	CreatedAt   time.Time
	UpdatedAt   time.Time
}
