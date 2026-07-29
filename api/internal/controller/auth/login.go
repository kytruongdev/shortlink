package auth

import (
	"context"
	"errors"

	"github.com/kytruongdev/shortlink/internal/pkg/apperror"
	pkgauth "github.com/kytruongdev/shortlink/internal/pkg/auth"
	repouser "github.com/kytruongdev/shortlink/internal/repository/user"
)

// Login verifies credentials and returns freshly issued tokens.
func (i *impl) Login(ctx context.Context, email, password string) (AuthResult, error) {
	email = normalizeEmail(email)

	u, err := i.users.GetByEmail(ctx, email)
	if err != nil {
		if errors.Is(err, repouser.ErrNotFound) {
			// Same response as a wrong password so an attacker cannot tell which emails exist.
			return AuthResult{}, apperror.Unauthorized(codeInvalidCredentials, "invalid email or password")
		}
		return AuthResult{}, err
	}

	if err := pkgauth.ComparePassword(u.PasswordHash, password); err != nil {
		return AuthResult{}, apperror.Unauthorized(codeInvalidCredentials, "invalid email or password")
	}

	res, err := i.issueTokens(ctx, u)
	if err != nil {
		return AuthResult{}, err
	}
	res.User = u
	return res, nil
}
