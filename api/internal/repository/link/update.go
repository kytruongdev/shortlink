package link

import (
	"context"
	"errors"

	"github.com/google/uuid"
	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgconn"
	pkgerrors "github.com/pkg/errors"

	"github.com/kytruongdev/shortlink/internal/model"
	"github.com/kytruongdev/shortlink/internal/repository/sqlc"
)

// UpdateURL changes the destination of the user's link. It returns ErrNotFound
// when the link isn't the user's, and ErrConflict when the new URL already exists.
func (i *impl) UpdateURL(ctx context.Context, code string, userID uuid.UUID, original, normalized string) (model.Link, error) {
	row, err := i.q.UpdateLinkURL(ctx, sqlc.UpdateLinkURLParams{
		Code:          code,
		OriginalUrl:   original,
		NormalizedUrl: normalized,
		UserID:        &userID,
	})
	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return model.Link{}, ErrNotFound
		}
		var pgErr *pgconn.PgError
		if errors.As(err, &pgErr) && pgErr.Code == uniqueViolationCode {
			return model.Link{}, ErrConflict
		}
		return model.Link{}, pkgerrors.WithStack(err)
	}
	return toModel(row), nil
}
