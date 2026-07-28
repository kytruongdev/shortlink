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

func TestListByUserID(t *testing.T) {
	tcs := map[string]struct {
		userID    uuid.UUID
		wantCodes []string
	}{
		"returns the user's links newest first": {
			userID:    uuid.MustParse(testfixture.SeededUserID),
			wantCodes: []string{"newer1", "older1"},
		},
		"returns empty for a user with no links": {
			userID:    uuid.New(),
			wantCodes: []string{},
		},
	}

	for name, tc := range tcs {
		t.Run(name, func(t *testing.T) {
			testutil.WithTxDB(t, func(tx pgx.Tx) {
				testutil.LoadSQLFile(t, tx, "testfixture/user_links.sql")

				got, err := New(tx).ListByUserID(context.Background(), tc.userID)
				require.NoError(t, err)

				codes := make([]string, len(got))
				for i, l := range got {
					codes[i] = l.Code
				}
				assert.Equal(t, tc.wantCodes, codes)
			})
		})
	}
}
