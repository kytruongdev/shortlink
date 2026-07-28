package user

import (
	"context"
	"testing"

	"github.com/google/uuid"
	"github.com/jackc/pgx/v5"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"

	"github.com/kytruongdev/shortlink/internal/model"
	"github.com/kytruongdev/shortlink/internal/pkg/testutil"
	"github.com/kytruongdev/shortlink/internal/repository/user/testfixture"
)

func TestCreate(t *testing.T) {
	tcs := map[string]struct {
		fixture string
		given   model.User
		wantErr error
	}{
		"inserts user and sets id and created_at": {
			given: model.User{Email: "new@example.com", PasswordHash: "hash"},
		},
		"duplicate email returns ErrConflict": {
			fixture: "testfixture/users.sql",
			given:   model.User{Email: testfixture.SeededEmail, PasswordHash: "hash"},
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
				assert.Equal(t, tc.given.Email, got.Email)
				assert.NotEqual(t, uuid.Nil, got.ID)
				assert.False(t, got.CreatedAt.IsZero())
			})
		})
	}
}
