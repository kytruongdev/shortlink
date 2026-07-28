package auth

import (
	"context"

	pkgauth "github.com/kytruongdev/shortlink/internal/pkg/auth"
)

// Logout revokes the presented refresh token; it is idempotent and never errors on an unknown token.
func (i *impl) Logout(ctx context.Context, refreshToken string) error {
	if refreshToken == "" {
		return nil
	}
	return i.tokens.Revoke(ctx, pkgauth.HashRefreshToken(refreshToken))
}
