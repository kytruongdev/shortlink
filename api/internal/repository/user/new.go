package user

import (
	"context"

	"github.com/google/uuid"

	"github.com/kytruongdev/shortlink/internal/model"
	"github.com/kytruongdev/shortlink/internal/repository/sqlc"
)

// Repository manages data access for users.
type Repository interface {
	Create(ctx context.Context, user model.User) (model.User, error)
	GetByEmail(ctx context.Context, email string) (model.User, error)
	GetByID(ctx context.Context, id uuid.UUID) (model.User, error)
}

type impl struct {
	q *sqlc.Queries
}

// New returns a Repository backed by the given database handle (pool or tx).
func New(db sqlc.DBTX) Repository {
	return &impl{q: sqlc.New(db)}
}

func toModel(u sqlc.User) model.User {
	return model.User{
		ID:           u.ID,
		Email:        u.Email,
		PasswordHash: u.PasswordHash,
		CreatedAt:    u.CreatedAt,
	}
}
