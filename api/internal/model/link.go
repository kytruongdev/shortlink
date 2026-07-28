package model

import (
	"time"

	"github.com/google/uuid"
)

const (
	// ShortCodeLength is the number of base62 characters in a generated code.
	ShortCodeLength = 7
	// MaxURLLen is the maximum accepted length of an input URL.
	MaxURLLen = 2048
)

// Link is a shortened-URL mapping.
type Link struct {
	Code          string
	OriginalURL   string
	NormalizedURL string
	CreatedAt     time.Time
	UserID        *uuid.UUID // nil for anonymous links
	CreatorIP     *string    // set for anonymous links, for per-IP quota
}
