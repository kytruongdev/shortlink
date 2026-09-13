package link

import (
	"context"

	"github.com/google/uuid"
	pkgerrors "github.com/pkg/errors"

	"github.com/kytruongdev/shortlink/internal/model"
)

// ListByUserID returns the user's links, newest first.
func (i *impl) ListByUserID(ctx context.Context, userID uuid.UUID) ([]model.Link, error) {
	rows, err := i.q.ListLinksByUserID(ctx, &userID)
	if err != nil {
		return nil, pkgerrors.WithStack(err)
	}

	links := make([]model.Link, 0, len(rows))
	for _, row := range rows {
		links = append(links, toModel(row))
	}
	return links, nil
}
