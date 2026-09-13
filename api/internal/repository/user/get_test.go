package user

import (
	"context"
	"testing"

	"github.com/google/uuid"
	"github.com/jackc/pgx/v5"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"

	"github.com/kytruongdev/shortlink/internal/pkg/testutil"
	"github.com/kytruongdev/shortlink/internal/repository/user/testfixture"
)

func TestGetByEmail(t *testing.T) {
	tcs := map[string]struct {
		fixture string
		email   string
		wantErr error
	}{
		"returns user for existing email": {
			fixture: "testfixture/users.sql",
			email:   testfixture.SeededEmail,
		},
		"returns ErrNotFound for missing email": {
			email:   "missing@example.com",
			wantErr: ErrNotFound,
		},
	}

	for name, tc := range tcs {
		t.Run(name, func(t *testing.T) {
			testutil.WithTxDB(t, func(tx pgx.Tx) {
				if tc.fixture != "" {
					testutil.LoadSQLFile(t, tx, tc.fixture)
				}

				got, err := New(tx).GetByEmail(context.Background(), tc.email)
				if tc.wantErr != nil {
					require.ErrorIs(t, err, tc.wantErr)
					return
				}
				require.NoError(t, err)
				assert.Equal(t, testfixture.SeededEmail, got.Email)
			})
		})
	}
}

func TestGetByID(t *testing.T) {
	tcs := map[string]struct {
		fixture string
		id      uuid.UUID
		wantErr error
	}{
		"returns user for existing id": {
			fixture: "testfixture/users.sql",
			id:      uuid.MustParse(testfixture.SeededUserID),
		},
		"returns ErrNotFound for missing id": {
			id:      uuid.New(),
			wantErr: ErrNotFound,
		},
	}

	for name, tc := range tcs {
		t.Run(name, func(t *testing.T) {
			testutil.WithTxDB(t, func(tx pgx.Tx) {
				if tc.fixture != "" {
					testutil.LoadSQLFile(t, tx, tc.fixture)
				}

				got, err := New(tx).GetByID(context.Background(), tc.id)
				if tc.wantErr != nil {
					require.ErrorIs(t, err, tc.wantErr)
					return
				}
				require.NoError(t, err)
				assert.Equal(t, testfixture.SeededEmail, got.Email)
			})
		})
	}
}
