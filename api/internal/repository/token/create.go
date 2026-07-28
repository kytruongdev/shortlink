package token

import (
	"context"

	pkgerrors "github.com/pkg/errors"

	"github.com/kytruongdev/shortlink/internal/model"
	"github.com/kytruongdev/shortlink/internal/repository/sqlc"
)

// Create inserts a refresh token record.
func (i *impl) Create(ctx context.Context, token model.RefreshToken) (model.RefreshToken, error) {
	row, err := i.q.CreateRefreshToken(ctx, sqlc.CreateRefreshTokenParams{
		UserID:    token.UserID,
		TokenHash: token.TokenHash,
		ExpiresAt: token.ExpiresAt,
	})
	if err != nil {
		return model.RefreshToken{}, pkgerrors.WithStack(err)
	}
	return toModel(row), nil
}
