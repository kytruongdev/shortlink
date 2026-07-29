package auth

import (
	"context"
	"errors"

	pkgerrors "github.com/pkg/errors"

	"github.com/kytruongdev/shortlink/internal/model"
	"github.com/kytruongdev/shortlink/internal/pkg/apperror"
	pkgauth "github.com/kytruongdev/shortlink/internal/pkg/auth"
	repouser "github.com/kytruongdev/shortlink/internal/repository/user"
)

// Register creates an account and returns freshly issued tokens (auto-login).
// The handler validates the input form; this normalizes the email for storage.
func (i *impl) Register(ctx context.Context, email, password string) (AuthResult, error) {
	email = normalizeEmail(email)

	hash, err := pkgauth.HashPassword(password)
	if err != nil {
		return AuthResult{}, pkgerrors.WithStack(err)
	}

	created, err := i.users.Create(ctx, model.User{Email: email, PasswordHash: hash})
	if err != nil {
		if errors.Is(err, repouser.ErrConflict) {
			return AuthResult{}, apperror.Conflict(codeEmailTaken, "email already registered")
		}
		return AuthResult{}, err
	}

	res, err := i.issueTokens(ctx, created)
	if err != nil {
		return AuthResult{}, err
	}
	res.User = created
	return res, nil
}
