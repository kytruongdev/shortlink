package link

import (
	"context"
	"testing"

	"github.com/jackc/pgx/v5"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"

	"github.com/kytruongdev/shortlink/internal/model"
	"github.com/kytruongdev/shortlink/internal/pkg/testutil"
	"github.com/kytruongdev/shortlink/internal/repository/link/testfixture"
)

func TestCreate(t *testing.T) {
	tcs := map[string]struct {
		fixture string
		given   model.Link
		wantErr error
	}{
		"inserts link and sets created_at": {
			given: model.Link{Code: "newcode", OriginalURL: "https://x.com/a", NormalizedURL: "https://x.com/a"},
		},
		"duplicate code returns ErrConflict": {
			fixture: "testfixture/links.sql",
			given:   model.Link{Code: testfixture.SeededCode, OriginalURL: "https://y.com", NormalizedURL: "https://y.com"},
			wantErr: ErrConflict,
		},
		"duplicate normalized url returns ErrConflict": {
			fixture: "testfixture/links.sql",
			given:   model.Link{Code: "another", OriginalURL: testfixture.SeededOriginalURL, NormalizedURL: testfixture.SeededNormalizedURL},
			wantErr: ErrConflict,
		},
	}

	for name, tc := range tcs {
		t.Run(name, func(t *testing.T) {
			testutil.WithTxDB(t, func(tx pgx.Tx) {
				if tc.fixture != "" {
					testutil.LoadSQLFile(t, tx, tc.fixture)
				}

				got, err := New(tx).Create(context.Background(), tc.given)
				if tc.wantErr != nil {
					require.ErrorIs(t, err, tc.wantErr)
					return
				}
				require.NoError(t, err)
				assert.Equal(t, tc.given.Code, got.Code)
				assert.False(t, got.CreatedAt.IsZero())
			})
		})
	}
}
