package model

import (
	"time"

	"github.com/google/uuid"
)

// User is a registered account.
type User struct {
	ID           uuid.UUID
	Email        string
	PasswordHash string
	CreatedAt    time.Time
}
