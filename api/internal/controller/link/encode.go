package link

import (
	"context"
	"errors"
	"log/slog"
	"time"

	"github.com/google/uuid"

	"github.com/kytruongdev/shortlink/internal/model"
	"github.com/kytruongdev/shortlink/internal/pkg/apperror"
	"github.com/kytruongdev/shortlink/internal/pkg/urlshortener"
	repolink "github.com/kytruongdev/shortlink/internal/repository/link"
)

// maxCodeAttempts caps code regeneration on collision.
const maxCodeAttempts = 3

// Encode returns the existing link for a known URL, otherwise generates a unique code and stores it.
// A logged-in userID owns the link and skips the quota; an anonymous caller is capped per IP per day.
func (i *impl) Encode(ctx context.Context, rawURL string, userID *uuid.UUID, clientIP string) (model.Link, error) {
	normalized, err := urlshortener.Normalize(rawURL)
	if err != nil {
		return model.Link{}, err
	}

	// Dedup first: a known URL returns its existing code and never counts against quota.
	existing, err := i.repo.GetByNormalizedURL(ctx, normalized)
	if err == nil {
		return existing, nil
	}
	if !errors.Is(err, repolink.ErrNotFound) {
		return model.Link{}, err
	}

	var creatorIP *string
	if userID == nil {
		n, countErr := i.repo.CountByCreatorIPSince(ctx, clientIP, startOfDayUTC())
		if countErr != nil {
			return model.Link{}, countErr
		}
		if n >= model.MaxAnonLinksPerDay {
			return model.Link{}, apperror.TooManyRequests(codeQuotaExceeded, "daily link limit reached")
		}
		creatorIP = &clientIP
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
			UserID:        userID,
			CreatorIP:     creatorIP,
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

// startOfDayUTC returns today's midnight in UTC, the window start for the per-IP daily quota.
func startOfDayUTC() time.Time {
	now := time.Now().UTC()
	return time.Date(now.Year(), now.Month(), now.Day(), 0, 0, 0, 0, time.UTC)
}
