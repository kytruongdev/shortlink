package link

import (
	"context"

	"github.com/google/uuid"
	pkgerrors "github.com/pkg/errors"

	"github.com/kytruongdev/shortlink/internal/repository/sqlc"
)

// Delete removes the user's link by code, returning ErrNotFound when nothing was deleted.
func (i *impl) Delete(ctx context.Context, code string, userID uuid.UUID) error {
	n, err := i.q.DeleteLink(ctx, sqlc.DeleteLinkParams{Code: code, UserID: &userID})
	if err != nil {
		return pkgerrors.WithStack(err)
	}
	if n == 0 {
		return ErrNotFound
	}
	return nil
}
