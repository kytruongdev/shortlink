package link

import (
	"context"
	"testing"
	"time"

	"github.com/jackc/pgx/v5"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"

	"github.com/kytruongdev/shortlink/internal/model"
	"github.com/kytruongdev/shortlink/internal/pkg/testutil"
)

func TestCountByCreatorIPSince(t *testing.T) {
	const ip = "203.0.113.5"

	tcs := map[string]struct {
		ip    string
		since time.Time
		want  int
	}{
		"counts the ip's links since a day ago": {ip: ip, since: time.Now().Add(-24 * time.Hour), want: 2},
		"counts nothing since a future time":    {ip: ip, since: time.Now().Add(time.Hour), want: 0},
		"counts nothing for a different ip":     {ip: "198.51.100.9", since: time.Now().Add(-24 * time.Hour), want: 0},
	}

	for name, tc := range tcs {
		t.Run(name, func(t *testing.T) {
			testutil.WithTxDB(t, func(tx pgx.Tx) {
				repo := New(tx)
				ipVal := ip
				for _, code := range []string{"cnt1", "cnt2"} {
					_, err := repo.Create(context.Background(), model.Link{
						Code: code, OriginalURL: "https://" + code + ".com", NormalizedURL: "https://" + code + ".com", CreatorIP: &ipVal,
					})
					require.NoError(t, err)
				}

				got, err := repo.CountByCreatorIPSince(context.Background(), tc.ip, tc.since)
				require.NoError(t, err)
				assert.Equal(t, tc.want, got)
			})
		})
	}
}
