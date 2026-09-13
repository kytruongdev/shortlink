package token

import (
	"context"

	pkgerrors "github.com/pkg/errors"
)

// Revoke marks the refresh token with the given hash as revoked; it is idempotent.
func (i *impl) Revoke(ctx context.Context, tokenHash string) error {
	if err := i.q.RevokeRefreshToken(ctx, tokenHash); err != nil {
		return pkgerrors.WithStack(err)
	}
	return nil
}
