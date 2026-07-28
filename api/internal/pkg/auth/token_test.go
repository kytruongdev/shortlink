package auth

import (
	"testing"
	"time"

	"github.com/golang-jwt/jwt/v5"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestAccessToken(t *testing.T) {
	secret := []byte("0123456789abcdef0123456789abcdef")
	const userID = "user-123"

	tcs := map[string]struct {
		token       func(t *testing.T) string
		parseSecret []byte
		wantErr     bool
	}{
		"sign then parse round-trips the subject": {
			token: func(t *testing.T) string {
				raw, err := SignAccessToken(secret, userID, time.Minute)
				require.NoError(t, err)
				return raw
			},
			parseSecret: secret,
		},
		"expired token rejected": {
			token: func(t *testing.T) string {
				raw, err := SignAccessToken(secret, userID, -time.Minute)
				require.NoError(t, err)
				return raw
			},
			parseSecret: secret,
			wantErr:     true,
		},
		"wrong secret rejected": {
			token: func(t *testing.T) string {
				raw, err := SignAccessToken(secret, userID, time.Minute)
				require.NoError(t, err)
				return raw
			},
			parseSecret: []byte("different-secret-different-secret"),
			wantErr:     true,
		},
		"alg none rejected (algorithm confusion)": {
			token: func(t *testing.T) string {
				tok := jwt.NewWithClaims(jwt.SigningMethodNone, jwt.RegisteredClaims{
					Subject:   userID,
					ExpiresAt: jwt.NewNumericDate(time.Now().Add(time.Minute)),
				})
				raw, err := tok.SignedString(jwt.UnsafeAllowNoneSignatureType)
				require.NoError(t, err)
				return raw
			},
			parseSecret: secret,
			wantErr:     true,
		},
	}

	for name, tc := range tcs {
		t.Run(name, func(t *testing.T) {
			claims, err := ParseAccessToken(tc.parseSecret, tc.token(t))
			if tc.wantErr {
				require.Error(t, err)
				return
			}
			require.NoError(t, err)
			assert.Equal(t, userID, claims.UserID)
			assert.WithinDuration(t, time.Now().Add(time.Minute), claims.ExpiresAt, 5*time.Second)
		})
	}
}
