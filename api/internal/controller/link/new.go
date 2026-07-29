package link

import (
	"context"

	"github.com/google/uuid"

	"github.com/kytruongdev/shortlink/internal/model"
	repolink "github.com/kytruongdev/shortlink/internal/repository/link"
)

// Controller is the link use-case interface.
type Controller interface {
	Encode(ctx context.Context, rawURL string, userID *uuid.UUID, clientIP string) (model.Link, error)
	Decode(ctx context.Context, code string) (model.Link, error)
	ListByUser(ctx context.Context, userID uuid.UUID) ([]model.Link, error)
}

type impl struct {
	repo repolink.Repository
}

// New creates a Controller backed by the given repository.
func New(repo repolink.Repository) Controller {
	return &impl{repo: repo}
}
