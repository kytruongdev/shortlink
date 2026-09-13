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

func TestRefresh(t *testing.T) {
	userID := uuid.New()
	const raw = "raw-refresh-token"
	hash := pkgauth.HashRefreshToken(raw)
	revokedAt := time.Now()

	tcs := map[string]struct {
		token      string
		setup      func(u *repouser.MockRepository, tk *repotoken.MockRepository)
		wantStatus int // 0 = success
	}{
		"rotates a valid token": {
			token: raw,
			setup: func(u *repouser.MockRepository, tk *repotoken.MockRepository) {
				tk.EXPECT().GetByHash(mock.Anything, hash).
					Return(model.RefreshToken{UserID: userID, TokenHash: hash, ExpiresAt: time.Now().Add(time.Hour)}, nil).Once()
				tk.EXPECT().Revoke(mock.Anything, hash).Return(nil).Once()
				u.EXPECT().GetByID(mock.Anything, userID).Return(model.User{ID: userID, Email: "user@example.com"}, nil).Once()
				tk.EXPECT().Create(mock.Anything, mock.Anything).Return(model.RefreshToken{}, nil).Once()
			},
		},
		"revoked token is unauthorized": {
			token: raw,
			setup: func(_ *repouser.MockRepository, tk *repotoken.MockRepository) {
				tk.EXPECT().GetByHash(mock.Anything, hash).
					Return(model.RefreshToken{UserID: userID, ExpiresAt: time.Now().Add(time.Hour), RevokedAt: &revokedAt}, nil).Once()
			},
			wantStatus: http.StatusUnauthorized,
		},
		"expired token is unauthorized": {
			token: raw,
			setup: func(_ *repouser.MockRepository, tk *repotoken.MockRepository) {
				tk.EXPECT().GetByHash(mock.Anything, hash).
					Return(model.RefreshToken{UserID: userID, ExpiresAt: time.Now().Add(-time.Hour)}, nil).Once()
			},
			wantStatus: http.StatusUnauthorized,
		},
		"unknown token is unauthorized": {
			token: raw,
			setup: func(_ *repouser.MockRepository, tk *repotoken.MockRepository) {
				tk.EXPECT().GetByHash(mock.Anything, hash).Return(model.RefreshToken{}, repotoken.ErrNotFound).Once()
			},
			wantStatus: http.StatusUnauthorized,
		},
		"empty token is unauthorized": {
			token:      "",
			setup:      func(*repouser.MockRepository, *repotoken.MockRepository) {},
			wantStatus: http.StatusUnauthorized,
		},
	}

	for name, tc := range tcs {
		t.Run(name, func(t *testing.T) {
			users := repouser.NewMockRepository(t)
			tokens := repotoken.NewMockRepository(t)
			tc.setup(users, tokens)

			ctrl := New(users, tokens, testSecret, testAccessTTL, testRefreshTTL)
			res, err := ctrl.Refresh(context.Background(), tc.token)

			if tc.wantStatus != 0 {
				var appErr *apperror.Error
				require.ErrorAs(t, err, &appErr)
				assert.Equal(t, tc.wantStatus, appErr.HTTPStatus())
				return
			}
			require.NoError(t, err)
			assert.NotEmpty(t, res.AccessToken)
			assert.NotEmpty(t, res.RefreshToken)

			claims, err := pkgauth.ParseAccessToken(testSecret, res.AccessToken)
			require.NoError(t, err)
			assert.Equal(t, userID.String(), claims.UserID)
			assert.Equal(t, "user", claims.Username)
		})
	}
}
