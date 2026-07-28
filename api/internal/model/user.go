package model

import (
	"time"

	"github.com/google/uuid"
)

const (
	// MinPasswordLen is the minimum accepted password length.
	MinPasswordLen = 8
	// MaxPasswordLen is the maximum accepted password length.
	MaxPasswordLen = 24
)

// User is a registered account.
type User struct {
	ID           uuid.UUID
	Email        string
	PasswordHash string
	CreatedAt    time.Time
}
