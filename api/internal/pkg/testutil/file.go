package testutil

import (
	"context"
	"os"
	"testing"

	"github.com/jackc/pgx/v5"
	"github.com/stretchr/testify/require"
)

// LoadSQLFile executes the SQL file at path inside tx; path is relative to the test's working directory.
func LoadSQLFile(t *testing.T, tx pgx.Tx, path string) {
	t.Helper()
	b, err := os.ReadFile(path)
	require.NoError(t, err)
	_, err = tx.Exec(context.Background(), string(b))
	require.NoError(t, err)
}
