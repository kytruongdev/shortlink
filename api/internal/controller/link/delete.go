package link

import (
	"context"
	"errors"

	"github.com/google/uuid"

	"github.com/kytruongdev/shortlink/internal/pkg/apperror"
	repolink "github.com/kytruongdev/shortlink/internal/repository/link"
)

// Delete removes the user's link, or returns a not-found error.
func (i *impl) Delete(ctx context.Context, code string, userID uuid.UUID) error {
	if err := i.repo.Delete(ctx, code, userID); err != nil {
		if errors.Is(err, repolink.ErrNotFound) {
			return apperror.NotFound(codeNotFound, "link not found")
		}
		return err
	}
	return nil
}
