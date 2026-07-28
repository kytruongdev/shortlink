package token

import (
	"context"

	"github.com/kytruongdev/shortlink/internal/model"
	"github.com/kytruongdev/shortlink/internal/repository/sqlc"
)

// Repository manages data access for refresh tokens.
type Repository interface {
	Create(ctx context.Context, token model.RefreshToken) (model.RefreshToken, error)
	GetByHash(ctx context.Context, tokenHash string) (model.RefreshToken, error)
	Revoke(ctx context.Context, tokenHash string) error
}

type impl struct {
	q *sqlc.Queries
}

// New returns a Repository backed by the given database handle (pool or tx).
func New(db sqlc.DBTX) Repository {
	return &impl{q: sqlc.New(db)}
}

func toModel(t sqlc.RefreshToken) model.RefreshToken {
	return model.RefreshToken{
		ID:        t.ID,
		UserID:    t.UserID,
		TokenHash: t.TokenHash,
		ExpiresAt: t.ExpiresAt,
		RevokedAt: t.RevokedAt,
		CreatedAt: t.CreatedAt,
	}
}
