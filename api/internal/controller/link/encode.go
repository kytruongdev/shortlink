package link

import (
	"context"
	"errors"
	"log/slog"

	"github.com/kytruongdev/shortlink/internal/model"
	"github.com/kytruongdev/shortlink/internal/pkg/apperror"
	"github.com/kytruongdev/shortlink/internal/pkg/urlshortener"
	repolink "github.com/kytruongdev/shortlink/internal/repository/link"
)

// maxCodeAttempts caps code regeneration on collision.
const maxCodeAttempts = 3

// Encode returns the existing link for a known URL, otherwise generates a unique code and stores it.
func (i *impl) Encode(ctx context.Context, rawURL string) (model.Link, error) {
	normalized, err := urlshortener.Normalize(rawURL)
	if err != nil {
		return model.Link{}, err
	}

	existing, err := i.repo.GetByNormalizedURL(ctx, normalized)
	if err == nil {
		return existing, nil
	}
	if !errors.Is(err, repolink.ErrNotFound) {
		return model.Link{}, err
	}

	for attempt := range maxCodeAttempts {
		code, genErr := urlshortener.Generate(model.ShortCodeLength)
		if genErr != nil {
			return model.Link{}, genErr
		}

		link, createErr := i.repo.Create(ctx, model.Link{
			Code:          code,
			OriginalURL:   rawURL,
			NormalizedURL: normalized,
		})
		if createErr == nil {
			return link, nil
		}
		if !errors.Is(createErr, repolink.ErrConflict) {
			return model.Link{}, createErr
		}

		// Conflict means a code collision (retry) or a concurrent insert of the same URL (return it).
		if found, getErr := i.repo.GetByNormalizedURL(ctx, normalized); getErr == nil {
			return found, nil
		}
		slog.Warn("short code collision, retrying", slog.String("code", code), slog.Int("attempt", attempt+1))
	}

	slog.Error("failed to generate unique code", slog.Int("attempts", maxCodeAttempts))
	return model.Link{}, apperror.Internal(codeGenerationFailed, "failed to generate unique code")
}
