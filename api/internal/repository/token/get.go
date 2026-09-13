package token

import (
	"context"
	"errors"

	"github.com/jackc/pgx/v5"
	pkgerrors "github.com/pkg/errors"

	"github.com/kytruongdev/shortlink/internal/model"
)

// GetByHash returns the refresh token with the given hash, or ErrNotFound.
func (i *impl) GetByHash(ctx context.Context, tokenHash string) (model.RefreshToken, error) {
	row, err := i.q.GetRefreshTokenByHash(ctx, tokenHash)
	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return model.RefreshToken{}, ErrNotFound
		}
		return model.RefreshToken{}, pkgerrors.WithStack(err)
	}
	return toModel(row), nil
}
