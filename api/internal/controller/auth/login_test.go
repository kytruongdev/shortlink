package auth

import (
	"context"
	"net/http"
	"testing"

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

func TestLogin(t *testing.T) {
	const (
		email    = "user@example.com"
		password = "password123"
	)
	hash, err := pkgauth.HashPassword(password)
	require.NoError(t, err)
	stored := model.User{ID: uuid.New(), Email: email, PasswordHash: hash}

	tcs := map[string]struct {
		password   string
		setup      func(u *repouser.MockRepository, tk *repotoken.MockRepository)
		wantStatus int // 0 = success
		wantCode   string
	}{
		"logs in with correct credentials": {
			password: password,
			setup: func(u *repouser.MockRepository, tk *repotoken.MockRepository) {
				u.EXPECT().GetByEmail(mock.Anything, email).Return(stored, nil).Once()
				tk.EXPECT().Create(mock.Anything, mock.Anything).Return(model.RefreshToken{}, nil).Once()
			},
		},
		"wrong password is unauthorized": {
			password: "wrong-password",
			setup: func(u *repouser.MockRepository, _ *repotoken.MockRepository) {
				u.EXPECT().GetByEmail(mock.Anything, email).Return(stored, nil).Once()
			},
			wantStatus: http.StatusUnauthorized,
			wantCode:   codeInvalidCredentials,
		},
		"unknown email is unauthorized with the same code (no enumeration)": {
			password: password,
			setup: func(u *repouser.MockRepository, _ *repotoken.MockRepository) {
				u.EXPECT().GetByEmail(mock.Anything, email).Return(model.User{}, repouser.ErrNotFound).Once()
			},
			wantStatus: http.StatusUnauthorized,
			wantCode:   codeInvalidCredentials,
		},
	}

	for name, tc := range tcs {
		t.Run(name, func(t *testing.T) {
			users := repouser.NewMockRepository(t)
			tokens := repotoken.NewMockRepository(t)
			tc.setup(users, tokens)

			ctrl := New(users, tokens, testSecret, testAccessTTL, testRefreshTTL)
			res, err := ctrl.Login(context.Background(), email, tc.password)

			if tc.wantStatus != 0 {
				var appErr *apperror.Error
				require.ErrorAs(t, err, &appErr)
				assert.Equal(t, tc.wantStatus, appErr.HTTPStatus())
				assert.Equal(t, tc.wantCode, appErr.Code)
				return
			}
			require.NoError(t, err)
			assert.Equal(t, stored, res.User)
			assert.NotEmpty(t, res.AccessToken)
			assert.NotEmpty(t, res.RefreshToken)
		})
	}
}
