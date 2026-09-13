package auth

import (
	"context"
	"time"

	"github.com/kytruongdev/shortlink/internal/model"
	repotoken "github.com/kytruongdev/shortlink/internal/repository/token"
	repouser "github.com/kytruongdev/shortlink/internal/repository/user"
)

// AuthResult is the outcome of a successful authentication.
type AuthResult struct {
	User         model.User
	AccessToken  string
	RefreshToken string        // raw token; the handler sets it as an httpOnly cookie
	ExpiresIn    time.Duration // access token lifetime
}

// Controller is the authentication use-case interface.
type Controller interface {
	Register(ctx context.Context, email, password string) (AuthResult, error)
	Login(ctx context.Context, email, password string) (AuthResult, error)
	Refresh(ctx context.Context, refreshToken string) (AuthResult, error)
	Logout(ctx context.Context, refreshToken string) error
}

type impl struct {
	users      repouser.Repository
	tokens     repotoken.Repository
	secret     []byte
	accessTTL  time.Duration
	refreshTTL time.Duration
}

// New creates a Controller backed by the user and refresh-token repositories.
func New(users repouser.Repository, tokens repotoken.Repository, secret []byte, accessTTL, refreshTTL time.Duration) Controller {
	return &impl{users: users, tokens: tokens, secret: secret, accessTTL: accessTTL, refreshTTL: refreshTTL}
}
