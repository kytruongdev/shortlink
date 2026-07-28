package auth

import (
	"context"
	"strings"
	"time"

	"github.com/google/uuid"
	pkgerrors "github.com/pkg/errors"

	"github.com/kytruongdev/shortlink/internal/model"
	pkgauth "github.com/kytruongdev/shortlink/internal/pkg/auth"
)

// issueTokens mints an access token and a rotated refresh token for userID,
// persisting only the refresh token's hash.
func (i *impl) issueTokens(ctx context.Context, userID uuid.UUID) (AuthResult, error) {
	access, err := pkgauth.SignAccessToken(i.secret, userID.String(), i.accessTTL)
	if err != nil {
		return AuthResult{}, pkgerrors.WithStack(err)
	}

	raw, hash, err := pkgauth.NewRefreshToken()
	if err != nil {
		return AuthResult{}, pkgerrors.WithStack(err)
	}

	if _, err := i.tokens.Create(ctx, model.RefreshToken{
		UserID:    userID,
		TokenHash: hash,
		ExpiresAt: time.Now().Add(i.refreshTTL),
	}); err != nil {
		return AuthResult{}, err
	}

	return AuthResult{AccessToken: access, RefreshToken: raw, ExpiresIn: i.accessTTL}, nil
}

// normalizeEmail lowercases and trims an email so lookups and uniqueness ignore case.
func normalizeEmail(email string) string {
	return strings.ToLower(strings.TrimSpace(email))
}
