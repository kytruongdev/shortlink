package testutil

import (
	"context"
	"os"
	"sync"
	"testing"

	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/stretchr/testify/require"
)

var (
	poolOnce sync.Once
	pool     *pgxpool.Pool
)

// WithTxDB runs callback inside a transaction that is rolled back afterwards, leaving the test database unchanged.
func WithTxDB(t *testing.T, callback func(pgx.Tx)) {
	t.Helper()
	poolOnce.Do(func() {
		var err error
		pool, err = pgxpool.New(context.Background(), os.Getenv("PG_URL"))
		require.NoError(t, err)
	})

	tx, err := pool.Begin(context.Background())
	require.NoError(t, err)
	defer tx.Rollback(context.Background()) //nolint:errcheck

	callback(tx)
}
