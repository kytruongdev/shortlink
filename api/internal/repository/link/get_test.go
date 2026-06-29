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

func TestGetByCode(t *testing.T) {
	tcs := map[string]struct {
		fixture string
		code    string
		wantURL string
		wantErr error
	}{
		"returns link for existing code": {
			fixture: "testfixture/links.sql",
			code:    testfixture.SeededCode,
			wantURL: testfixture.SeededOriginalURL,
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

				got, err := New(tx).GetByCode(context.Background(), tc.code)
				if tc.wantErr != nil {
					require.ErrorIs(t, err, tc.wantErr)
					return
				}
				require.NoError(t, err)
				assert.Equal(t, tc.wantURL, got.OriginalURL)
			})
		})
	}
}

func TestGetByNormalizedURL(t *testing.T) {
	tcs := map[string]struct {
		fixture       string
		normalizedURL string
		wantCode      string
		wantErr       error
	}{
		"returns link for existing url": {
			fixture:       "testfixture/links.sql",
			normalizedURL: testfixture.SeededNormalizedURL,
			wantCode:      testfixture.SeededCode,
		},
		"returns ErrNotFound for missing url": {
			normalizedURL: "https://nope.com",
			wantErr:       ErrNotFound,
		},
	}

	for name, tc := range tcs {
		t.Run(name, func(t *testing.T) {
			testutil.WithTxDB(t, func(tx pgx.Tx) {
				if tc.fixture != "" {
					testutil.LoadSQLFile(t, tx, tc.fixture)
				}

				got, err := New(tx).GetByNormalizedURL(context.Background(), tc.normalizedURL)
				if tc.wantErr != nil {
					require.ErrorIs(t, err, tc.wantErr)
					return
				}
				require.NoError(t, err)
				assert.Equal(t, tc.wantCode, got.Code)
			})
		})
	}
}
