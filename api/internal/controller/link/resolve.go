package link

import (
	"context"
	"errors"

	"github.com/kytruongdev/shortlink/internal/model"
	"github.com/kytruongdev/shortlink/internal/pkg/apperror"
	repolink "github.com/kytruongdev/shortlink/internal/repository/link"
)

// Resolve returns the link for code and counts the visit, or a not-found error.
func (i *impl) Resolve(ctx context.Context, code string) (model.Link, error) {
	link, err := i.repo.IncrementAndGetByCode(ctx, code)
	if err != nil {
		if errors.Is(err, repolink.ErrNotFound) {
			return model.Link{}, apperror.NotFound(codeNotFound, "short url not found")
		}
		return model.Link{}, err
	}
	return link, nil
}
