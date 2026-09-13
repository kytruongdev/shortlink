package auth

import (
	"context"
	"errors"
	"time"

	"github.com/kytruongdev/shortlink/internal/pkg/apperror"
	pkgauth "github.com/kytruongdev/shortlink/internal/pkg/auth"
	repotoken "github.com/kytruongdev/shortlink/internal/repository/token"
)

// Refresh validates a refresh token, rotates it, and returns a new token pair.
func (i *impl) Refresh(ctx context.Context, refreshToken string) (AuthResult, error) {
	invalid := apperror.Unauthorized(codeInvalidRefresh, "invalid refresh token")
	if refreshToken == "" {
		return AuthResult{}, invalid
	}

	hash := pkgauth.HashRefreshToken(refreshToken)
	stored, err := i.tokens.GetByHash(ctx, hash)
	if err != nil {
		if errors.Is(err, repotoken.ErrNotFound) {
			return AuthResult{}, invalid
		}
		return AuthResult{}, err
	}
	if stored.RevokedAt != nil || time.Now().After(stored.ExpiresAt) {
		return AuthResult{}, invalid
	}

	// Rotation: revoke the presented token before issuing its replacement.
	if err := i.tokens.Revoke(ctx, hash); err != nil {
		return AuthResult{}, err
	}

	// Reload the user so the reissued token carries the current username claim.
	user, err := i.users.GetByID(ctx, stored.UserID)
	if err != nil {
		return AuthResult{}, err
	}
	return i.issueTokens(ctx, user)
}
