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

func TestRevoke(t *testing.T) {
	testutil.WithTxDB(t, func(tx pgx.Tx) {
		testutil.LoadSQLFile(t, tx, "testfixture/refresh_tokens.sql")
		repo := New(tx)

		require.NoError(t, repo.Revoke(context.Background(), testfixture.SeededTokenHash))

		got, err := repo.GetByHash(context.Background(), testfixture.SeededTokenHash)
		require.NoError(t, err)
		assert.NotNil(t, got.RevokedAt, "revoke sets revoked_at")
	})
}
