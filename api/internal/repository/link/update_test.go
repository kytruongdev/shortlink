package link

import (
	"context"
	"testing"

	"github.com/google/uuid"
	"github.com/jackc/pgx/v5"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"

	"github.com/kytruongdev/shortlink/internal/pkg/testutil"
	"github.com/kytruongdev/shortlink/internal/repository/link/testfixture"
)

func TestUpdateURL(t *testing.T) {
	owner := uuid.MustParse(testfixture.SeededUserID)

	tcs := map[string]struct {
		code    string
		userID  uuid.UUID
		newURL  string
		newNorm string
		wantErr error
	}{
		"updates an owned link's destination": {
			code: "older1", userID: owner, newURL: "https://new.com", newNorm: "https://new.com",
		},
		"conflict with an existing url": {
			// newer1 already holds https://b.com
			code: "older1", userID: owner, newURL: "https://b.com", newNorm: "https://b.com", wantErr: ErrConflict,
		},
		"not owned returns ErrNotFound": {
			code: "newer1", userID: uuid.New(), newURL: "https://x.com", newNorm: "https://x.com", wantErr: ErrNotFound,
		},
	}

	for name, tc := range tcs {
		t.Run(name, func(t *testing.T) {
			testutil.WithTxDB(t, func(tx pgx.Tx) {
				testutil.LoadSQLFile(t, tx, "testfixture/user_links.sql")
				repo := New(tx)

				got, err := repo.UpdateURL(context.Background(), tc.code, tc.userID, tc.newURL, tc.newNorm)
				if tc.wantErr != nil {
					require.ErrorIs(t, err, tc.wantErr)
					return
				}
				require.NoError(t, err)
				assert.Equal(t, tc.newURL, got.OriginalURL)
				assert.Equal(t, tc.newNorm, got.NormalizedURL)
			})
		})
	}
}
