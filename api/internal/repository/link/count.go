package link

import (
	"context"
	"time"

	pkgerrors "github.com/pkg/errors"

	"github.com/kytruongdev/shortlink/internal/repository/sqlc"
)

// CountByCreatorIPSince counts links created by ip at or after the given time.
func (i *impl) CountByCreatorIPSince(ctx context.Context, ip string, since time.Time) (int, error) {
	n, err := i.q.CountLinksByCreatorIPSince(ctx, sqlc.CountLinksByCreatorIPSinceParams{
		CreatorIp: &ip,
		CreatedAt: since,
	})
	if err != nil {
		return 0, pkgerrors.WithStack(err)
	}
	return int(n), nil
}
