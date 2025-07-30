package domain

import (
	"github.com/google/uuid"
	"time"
)

type PasswordReset struct {
	Token     string
	UserID    uuid.UUID
	ExpiresAt time.Time
}
