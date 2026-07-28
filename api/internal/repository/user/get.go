package user

import (
	"context"
	"errors"

	"github.com/google/uuid"
	"github.com/jackc/pgx/v5"
	pkgerrors "github.com/pkg/errors"

	"github.com/kytruongdev/shortlink/internal/model"
)

// GetByEmail returns the user with the given email, or ErrNotFound.
func (i *impl) GetByEmail(ctx context.Context, email string) (model.User, error) {
	row, err := i.q.GetUserByEmail(ctx, email)
	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return model.User{}, ErrNotFound
		}
		return model.User{}, pkgerrors.WithStack(err)
	}
	return toModel(row), nil
}

// GetByID returns the user with the given id, or ErrNotFound.
func (i *impl) GetByID(ctx context.Context, id uuid.UUID) (model.User, error) {
	row, err := i.q.GetUserByID(ctx, id)
	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return model.User{}, ErrNotFound
		}
		return model.User{}, pkgerrors.WithStack(err)
	}
	return toModel(row), nil
}
