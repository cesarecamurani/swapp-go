package domain

import (
	"github.com/google/uuid"
)

type Item struct {
	ID          uuid.UUID
	Name        string
	Description string
	PictureURL  string
	UserID      uuid.UUID
}
