package link

import (
	"context"
	"errors"

	"github.com/jackc/pgx/v5/pgconn"
	pkgerrors "github.com/pkg/errors"

	"github.com/kytruongdev/shortlink/internal/model"
	"github.com/kytruongdev/shortlink/internal/repository/sqlc"
)

// Create inserts a link, returning ErrConflict on a unique-constraint violation.
func (i *impl) Create(ctx context.Context, link model.Link) (model.Link, error) {
	const uniqueViolationCode = "23505"

	row, err := i.q.CreateLink(ctx, sqlc.CreateLinkParams{
		Code:          link.Code,
		OriginalUrl:   link.OriginalURL,
		NormalizedUrl: link.NormalizedURL,
	})
	if err != nil {
		var pgErr *pgconn.PgError
		if errors.As(err, &pgErr) && pgErr.Code == uniqueViolationCode {
			return model.Link{}, ErrConflict
		}
		return model.Link{}, pkgerrors.WithStack(err)
	}
	return toModel(row), nil
}
