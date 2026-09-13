package link

import (
	"context"
	"errors"

	"github.com/google/uuid"

	"github.com/kytruongdev/shortlink/internal/model"
	"github.com/kytruongdev/shortlink/internal/pkg/apperror"
	"github.com/kytruongdev/shortlink/internal/pkg/urlshortener"
	repolink "github.com/kytruongdev/shortlink/internal/repository/link"
)

// UpdateURL changes a link's destination. Returns 404 when the link isn't the
// user's, and 409 when the new URL is already shortened by another link.
func (i *impl) UpdateURL(ctx context.Context, code string, userID uuid.UUID, rawURL string) (model.Link, error) {
	normalized, err := urlshortener.Normalize(rawURL)
	if err != nil {
		return model.Link{}, err
	}

	link, err := i.repo.UpdateURL(ctx, code, userID, rawURL, normalized)
	if err != nil {
		if errors.Is(err, repolink.ErrNotFound) {
			return model.Link{}, apperror.NotFound(codeNotFound, "link not found")
		}
		if errors.Is(err, repolink.ErrConflict) {
			return model.Link{}, apperror.Conflict(codeURLTaken, "this url is already shortened")
		}
		return model.Link{}, err
	}
	return link, nil
}
