package auth

import (
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestRefreshToken(t *testing.T) {
	raw, hash, err := NewRefreshToken()
	require.NoError(t, err)

	assert.NotEmpty(t, raw)
	assert.NotEqual(t, raw, hash, "stored hash must never equal the client token")
	assert.Len(t, hash, 64, "sha-256 hex is 64 chars")
	assert.Equal(t, hash, HashRefreshToken(raw), "hash is deterministic for the same raw token")

	raw2, hash2, err := NewRefreshToken()
	require.NoError(t, err)
	assert.NotEqual(t, raw, raw2, "each token is unique")
	assert.NotEqual(t, hash, hash2)
}
