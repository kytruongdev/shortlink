package token

import (
	"context"
	"testing"

	"github.com/jackc/pgx/v5"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"

	"github.com/kytruongdev/shortlink/internal/pkg/testutil"
	"github.com/kytruongdev/shortlink/internal/repository/token/testfixture"
)

func TestGetByHash(t *testing.T) {
	tcs := map[string]struct {
		fixture   string
		tokenHash string
		wantErr   error
	}{
		"returns token for existing hash": {
			fixture:   "testfixture/refresh_tokens.sql",
			tokenHash: testfixture.SeededTokenHash,
		},
		"returns ErrNotFound for missing hash": {
			tokenHash: "missing",
			wantErr:   ErrNotFound,
		},
	}

	for name, tc := range tcs {
		t.Run(name, func(t *testing.T) {
			testutil.WithTxDB(t, func(tx pgx.Tx) {
				if tc.fixture != "" {
					testutil.LoadSQLFile(t, tx, tc.fixture)
				}

				got, err := New(tx).GetByHash(context.Background(), tc.tokenHash)
				if tc.wantErr != nil {
					require.ErrorIs(t, err, tc.wantErr)
					return
				}
				require.NoError(t, err)
				assert.Equal(t, testfixture.SeededTokenHash, got.TokenHash)
				assert.Nil(t, got.RevokedAt)
			})
		})
	}
}
