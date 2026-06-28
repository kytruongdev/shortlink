package model

import "time"

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
}
