package auth

import (
	"context"
	"net/http"
	"testing"
	"time"

	"github.com/google/uuid"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
	"github.com/stretchr/testify/require"

	"github.com/kytruongdev/shortlink/internal/model"
	"github.com/kytruongdev/shortlink/internal/pkg/apperror"
	pkgauth "github.com/kytruongdev/shortlink/internal/pkg/auth"
	repotoken "github.com/kytruongdev/shortlink/internal/repository/token"
	repouser "github.com/kytruongdev/shortlink/internal/repository/user"
)

var testSecret = []byte("test-secret-test-secret-test-secret")

const (
	testAccessTTL  = 5 * time.Minute
	testRefreshTTL = 168 * time.Hour
)

func TestRegister(t *testing.T) {
	const password = "password123"
	createdID := uuid.New()
	created := model.User{ID: createdID, Email: "new@example.com"}

	tcs := map[string]struct {
		email      string
		password   string
		setup      func(u *repouser.MockRepository, tk *repotoken.MockRepository, storedHash *string)
		wantStatus int // 0 = success
	}{
		"registers, normalizes email, issues tokens": {
			email:    "  New@Example.com  ",
			password: password,
			setup: func(u *repouser.MockRepository, tk *repotoken.MockRepository, storedHash *string) {
				u.EXPECT().Create(mock.Anything, mock.MatchedBy(func(usr model.User) bool {
					return usr.Email == "new@example.com" && pkgauth.ComparePassword(usr.PasswordHash, password) == nil
				})).Return(created, nil).Once()
				tk.EXPECT().Create(mock.Anything, mock.Anything).
					Run(func(_ context.Context, token model.RefreshToken) { *storedHash = token.TokenHash }).
					Return(model.RefreshToken{}, nil).Once()
			},
		},
		"duplicate email is a conflict": {
			email:    "dupe@example.com",
			password: password,
			setup: func(u *repouser.MockRepository, _ *repotoken.MockRepository, _ *string) {
				u.EXPECT().Create(mock.Anything, mock.Anything).Return(model.User{}, repouser.ErrConflict).Once()
			},
			wantStatus: http.StatusConflict,
		},
	}

	for name, tc := range tcs {
		t.Run(name, func(t *testing.T) {
			users := repouser.NewMockRepository(t)
			tokens := repotoken.NewMockRepository(t)
			var storedHash string
			tc.setup(users, tokens, &storedHash)

			ctrl := New(users, tokens, testSecret, testAccessTTL, testRefreshTTL)
			res, err := ctrl.Register(context.Background(), tc.email, tc.password)

			if tc.wantStatus != 0 {
				var appErr *apperror.Error
				require.ErrorAs(t, err, &appErr)
				assert.Equal(t, tc.wantStatus, appErr.HTTPStatus())
				return
			}
			require.NoError(t, err)
			assert.Equal(t, created, res.User)
			assert.Equal(t, testAccessTTL, res.ExpiresIn)
			assert.NotEmpty(t, res.RefreshToken)

			// The stored value is the hash, never the raw token handed to the client.
			assert.NotEqual(t, res.RefreshToken, storedHash)
			assert.Equal(t, pkgauth.HashRefreshToken(res.RefreshToken), storedHash)

			// The access token encodes the new user's id and display username (email local-part).
			claims, err := pkgauth.ParseAccessToken(testSecret, res.AccessToken)
			require.NoError(t, err)
			assert.Equal(t, createdID.String(), claims.UserID)
			assert.Equal(t, "new", claims.Username)
		})
	}
}
