package auth

import (
	"errors"
	"time"

	"github.com/golang-jwt/jwt/v5"
)

// Claims carries the identity extracted from a verified access token.
type Claims struct {
	UserID    string
	ExpiresAt time.Time
}

// SignAccessToken issues an HS256 access token for userID valid for ttl.
func SignAccessToken(secret []byte, userID string, ttl time.Duration) (string, error) {
	now := time.Now()
	claims := jwt.RegisteredClaims{
		Subject:   userID,
		IssuedAt:  jwt.NewNumericDate(now),
		ExpiresAt: jwt.NewNumericDate(now.Add(ttl)),
	}
	return jwt.NewWithClaims(jwt.SigningMethodHS256, claims).SignedString(secret)
}

// ParseAccessToken verifies raw against secret and returns its claims.
// It rejects any non-HMAC signing method to prevent algorithm-confusion attacks.
func ParseAccessToken(secret []byte, raw string) (Claims, error) {
	var claims jwt.RegisteredClaims
	_, err := jwt.ParseWithClaims(raw, &claims, func(t *jwt.Token) (any, error) {
		if _, ok := t.Method.(*jwt.SigningMethodHMAC); !ok {
			return nil, errors.New("auth: unexpected signing method")
		}
		return secret, nil
	}, jwt.WithExpirationRequired())
	if err != nil {
		return Claims{}, err
	}
	return Claims{UserID: claims.Subject, ExpiresAt: claims.ExpiresAt.Time}, nil
}
