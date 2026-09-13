package link

import (
	"context"
	"testing"

	"github.com/jackc/pgx/v5"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"

	"github.com/kytruongdev/shortlink/internal/pkg/testutil"
	"github.com/kytruongdev/shortlink/internal/repository/link/testfixture"
)

func TestIncrementAndGetByCode(t *testing.T) {
	tcs := map[string]struct {
		fixture string
		code    string
		wantErr error
	}{
		"increments and returns the link": {
			fixture: "testfixture/links.sql",
			code:    testfixture.SeededCode,
		},
		"returns ErrNotFound for missing code": {
			code:    "missing",
			wantErr: ErrNotFound,
		},
	}

	for name, tc := range tcs {
		t.Run(name, func(t *testing.T) {
			testutil.WithTxDB(t, func(tx pgx.Tx) {
				if tc.fixture != "" {
					testutil.LoadSQLFile(t, tx, tc.fixture)
				}
				repo := New(tx)

				first, err := repo.IncrementAndGetByCode(context.Background(), tc.code)
				if tc.wantErr != nil {
					require.ErrorIs(t, err, tc.wantErr)
					return
				}
				require.NoError(t, err)
				assert.Equal(t, 1, first.ClickCount)

				second, err := repo.IncrementAndGetByCode(context.Background(), tc.code)
				require.NoError(t, err)
				assert.Equal(t, 2, second.ClickCount)
			})
		})
	}
}
