package model

import (
	"time"

	"github.com/google/uuid"
)

// RefreshToken is a persisted refresh token (stored as a hash) tied to a user session.
type RefreshToken struct {
	ID        uuid.UUID
	UserID    uuid.UUID
	TokenHash string
	ExpiresAt time.Time
	RevokedAt *time.Time
	CreatedAt time.Time
}
