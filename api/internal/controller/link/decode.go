package link

import (
	"context"
	"errors"

	"github.com/kytruongdev/shortlink/internal/model"
	"github.com/kytruongdev/shortlink/internal/pkg/apperror"
	repolink "github.com/kytruongdev/shortlink/internal/repository/link"
)

// Decode returns the link stored under code, or a not-found error.
func (i *impl) Decode(ctx context.Context, code string) (model.Link, error) {
	link, err := i.repo.GetByCode(ctx, code)
	if err != nil {
		if errors.Is(err, repolink.ErrNotFound) {
			return model.Link{}, apperror.NotFound(codeNotFound, "short url not found")
		}
		return model.Link{}, err
	}
	return link, nil
}
