package link

import (
	"context"
	"errors"

	"github.com/jackc/pgx/v5"
	pkgerrors "github.com/pkg/errors"

	"github.com/kytruongdev/shortlink/internal/model"
)

// GetByCode returns the link with the given code, or ErrNotFound.
func (i *impl) GetByCode(ctx context.Context, code string) (model.Link, error) {
	row, err := i.q.GetByCode(ctx, code)
	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return model.Link{}, ErrNotFound
		}
		return model.Link{}, pkgerrors.WithStack(err)
	}
	return toModel(row), nil
}

// GetByNormalizedURL returns the link with the given normalized URL, or ErrNotFound.
func (i *impl) GetByNormalizedURL(ctx context.Context, normalizedURL string) (model.Link, error) {
	row, err := i.q.GetByNormalizedURL(ctx, normalizedURL)
	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return model.Link{}, ErrNotFound
		}
		return model.Link{}, pkgerrors.WithStack(err)
	}
	return toModel(row), nil
}
