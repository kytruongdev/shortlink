package link

import (
	"context"
	"errors"

	"github.com/jackc/pgx/v5"
	pkgerrors "github.com/pkg/errors"

	"github.com/kytruongdev/shortlink/internal/model"
)

// IncrementAndGetByCode atomically bumps the click count and returns the link, or ErrNotFound.
func (i *impl) IncrementAndGetByCode(ctx context.Context, code string) (model.Link, error) {
	row, err := i.q.IncrementAndGetByCode(ctx, code)
	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return model.Link{}, ErrNotFound
		}
		return model.Link{}, pkgerrors.WithStack(err)
	}
	return toModel(row), nil
}
