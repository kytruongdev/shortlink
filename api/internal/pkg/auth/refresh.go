package auth

import (
	"crypto/rand"
	"crypto/sha256"
	"encoding/base64"
	"encoding/hex"
)

// refreshTokenBytes is the entropy of a raw refresh token before encoding.
const refreshTokenBytes = 32

// NewRefreshToken returns a new random refresh token and its SHA-256 hash.
// The raw token is handed to the client; only the hash is stored server-side.
func NewRefreshToken() (raw string, hash string, err error) {
	b := make([]byte, refreshTokenBytes)
	if _, err := rand.Read(b); err != nil {
		return "", "", err
	}
	raw = base64.RawURLEncoding.EncodeToString(b)
	return raw, HashRefreshToken(raw), nil
}

// HashRefreshToken returns the hex-encoded SHA-256 of a raw refresh token.
func HashRefreshToken(raw string) string {
	sum := sha256.Sum256([]byte(raw))
	return hex.EncodeToString(sum[:])
}
