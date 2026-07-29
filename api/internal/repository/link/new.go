package link

import (
	"context"
	"time"

	"github.com/google/uuid"

	"github.com/kytruongdev/shortlink/internal/model"
	"github.com/kytruongdev/shortlink/internal/repository/sqlc"
)

// Repository manages data access for links.
type Repository interface {
	Create(ctx context.Context, link model.Link) (model.Link, error)
	GetByCode(ctx context.Context, code string) (model.Link, error)
	GetByNormalizedURL(ctx context.Context, normalizedURL string) (model.Link, error)
	CountByCreatorIPSince(ctx context.Context, ip string, since time.Time) (int, error)
	ListByUserID(ctx context.Context, userID uuid.UUID) ([]model.Link, error)
	IncrementAndGetByCode(ctx context.Context, code string) (model.Link, error)
}

type impl struct {
	q *sqlc.Queries
}

// New returns a Repository backed by the given database handle (pool or tx).
func New(db sqlc.DBTX) Repository {
	return &impl{q: sqlc.New(db)}
}

func toModel(l sqlc.Link) model.Link {
	return model.Link{
		Code:          l.Code,
		OriginalURL:   l.OriginalUrl,
		NormalizedURL: l.NormalizedUrl,
		CreatedAt:     l.CreatedAt,
		UserID:        l.UserID,
		CreatorIP:     l.CreatorIp,
		ClickCount:    int(l.ClickCount),
	}
}
