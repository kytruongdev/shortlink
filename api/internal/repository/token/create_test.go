package token

import (
	"context"
	"testing"
	"time"

	"github.com/google/uuid"
	"github.com/jackc/pgx/v5"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"

	"github.com/kytruongdev/shortlink/internal/model"
	"github.com/kytruongdev/shortlink/internal/pkg/testutil"
	"github.com/kytruongdev/shortlink/internal/repository/token/testfixture"
)

func TestCreate(t *testing.T) {
	tcs := map[string]struct {
		userID    uuid.UUID
		wantError bool
	}{
		"inserts token for an existing user": {
			userID: uuid.MustParse(testfixture.SeededUserID),
		},
		"unknown user id violates the foreign key": {
			userID:    uuid.New(),
			wantError: true,
		},
	}

	for name, tc := range tcs {
		t.Run(name, func(t *testing.T) {
			testutil.WithTxDB(t, func(tx pgx.Tx) {
				testutil.LoadSQLFile(t, tx, "testfixture/refresh_tokens.sql")

				got, err := New(tx).Create(context.Background(), model.RefreshToken{
					UserID:    tc.userID,
					TokenHash: "freshhash",
					ExpiresAt: time.Now().Add(time.Hour),
				})
				if tc.wantError {
					require.Error(t, err)
					return
				}
				require.NoError(t, err)
				assert.NotEqual(t, uuid.Nil, got.ID)
				assert.Equal(t, "freshhash", got.TokenHash)
				assert.Nil(t, got.RevokedAt)
				assert.False(t, got.CreatedAt.IsZero())
			})
		})
	}
}
