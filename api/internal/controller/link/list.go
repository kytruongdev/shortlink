package link

import (
	"context"

	"github.com/google/uuid"

	"github.com/kytruongdev/shortlink/internal/model"
)

// ListByUser returns the links owned by the given user, newest first.
func (i *impl) ListByUser(ctx context.Context, userID uuid.UUID) ([]model.Link, error) {
	return i.repo.ListByUserID(ctx, userID)
}
