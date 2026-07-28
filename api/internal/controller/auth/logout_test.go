package auth

import (
	"context"
	"testing"

	"github.com/stretchr/testify/mock"
	"github.com/stretchr/testify/require"

	pkgauth "github.com/kytruongdev/shortlink/internal/pkg/auth"
	repotoken "github.com/kytruongdev/shortlink/internal/repository/token"
	repouser "github.com/kytruongdev/shortlink/internal/repository/user"
)

func TestLogout(t *testing.T) {
	const raw = "raw-refresh-token"

	tcs := map[string]struct {
		token string
		setup func(tk *repotoken.MockRepository)
	}{
		"revokes the presented token by its hash": {
			token: raw,
			setup: func(tk *repotoken.MockRepository) {
				tk.EXPECT().Revoke(mock.Anything, pkgauth.HashRefreshToken(raw)).Return(nil).Once()
			},
		},
		"empty token is a no-op": {
			token: "",
			setup: func(*repotoken.MockRepository) {},
		},
	}

	for name, tc := range tcs {
		t.Run(name, func(t *testing.T) {
			users := repouser.NewMockRepository(t)
			tokens := repotoken.NewMockRepository(t)
			tc.setup(tokens)

			ctrl := New(users, tokens, testSecret, testAccessTTL, testRefreshTTL)
			require.NoError(t, ctrl.Logout(context.Background(), tc.token))
		})
	}
}
