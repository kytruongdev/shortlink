package user

import (
	"context"
	"errors"

	"github.com/jackc/pgx/v5/pgconn"
	pkgerrors "github.com/pkg/errors"

	"github.com/kytruongdev/shortlink/internal/model"
	"github.com/kytruongdev/shortlink/internal/repository/sqlc"
)

// Create inserts a user, returning ErrConflict when the email already exists.
func (i *impl) Create(ctx context.Context, user model.User) (model.User, error) {
	const uniqueViolationCode = "23505"

	row, err := i.q.CreateUser(ctx, sqlc.CreateUserParams{
		Email:        user.Email,
		PasswordHash: user.PasswordHash,
	})
	if err != nil {
		var pgErr *pgconn.PgError
		if errors.As(err, &pgErr) && pgErr.Code == uniqueViolationCode {
			return model.User{}, ErrConflict
		}
		return model.User{}, pkgerrors.WithStack(err)
	}
	return toModel(row), nil
}
