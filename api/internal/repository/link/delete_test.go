package link

import (
	"context"
	"testing"

	"github.com/google/uuid"
	"github.com/jackc/pgx/v5"
	"github.com/stretchr/testify/require"

	"github.com/kytruongdev/shortlink/internal/pkg/testutil"
	"github.com/kytruongdev/shortlink/internal/repository/link/testfixture"
)

func TestDelete(t *testing.T) {
	owner := uuid.MustParse(testfixture.SeededUserID)

	tcs := map[string]struct {
		code    string
		userID  uuid.UUID
		wantErr error
	}{
		"deletes an owned link":            {code: "older1", userID: owner},
		"not owned returns ErrNotFound":    {code: "older1", userID: uuid.New(), wantErr: ErrNotFound},
		"missing code returns ErrNotFound": {code: "nope", userID: owner, wantErr: ErrNotFound},
	}

	for name, tc := range tcs {
		t.Run(name, func(t *testing.T) {
			testutil.WithTxDB(t, func(tx pgx.Tx) {
				testutil.LoadSQLFile(t, tx, "testfixture/user_links.sql")
				repo := New(tx)

				err := repo.Delete(context.Background(), tc.code, tc.userID)
				if tc.wantErr != nil {
					require.ErrorIs(t, err, tc.wantErr)
					return
				}
				require.NoError(t, err)

				_, err = repo.GetByCode(context.Background(), tc.code)
				require.ErrorIs(t, err, ErrNotFound)
			})
		})
	}
}
