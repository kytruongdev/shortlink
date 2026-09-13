package auth

import (
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestPassword(t *testing.T) {
	const password = "s3cret-password"

	hash, err := HashPassword(password)
	require.NoError(t, err)
	assert.NotEqual(t, password, hash, "plaintext must never be stored")

	tcs := map[string]struct {
		password string
		wantErr  bool
	}{
		"correct password matches": {password: password},
		"wrong password rejected":  {password: "wrong-password", wantErr: true},
		"empty password rejected":  {password: "", wantErr: true},
	}

	for name, tc := range tcs {
		t.Run(name, func(t *testing.T) {
			err := ComparePassword(hash, tc.password)
			if tc.wantErr {
				require.Error(t, err)
				return
			}
			require.NoError(t, err)
		})
	}
}
