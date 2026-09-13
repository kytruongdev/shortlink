package rest

import (
	"net/http"
	"net/http/httptest"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
	"github.com/stretchr/testify/require"

	ctrlauth "github.com/kytruongdev/shortlink/internal/controller/auth"
	"github.com/kytruongdev/shortlink/internal/infra/httpserver"
	"github.com/kytruongdev/shortlink/internal/pkg/apperror"
)

func TestRefreshHandler(t *testing.T) {
	result := ctrlauth.AuthResult{AccessToken: "new-access", RefreshToken: "new-refresh", ExpiresIn: 5 * time.Minute}

	tcs := map[string]struct {
		cookie     string // refresh cookie value; "" = no cookie sent
		setup      func(*ctrlauth.MockController)
		wantStatus int
		wantCookie bool
	}{
		"rotates the token with a valid cookie": {
			cookie: "old-refresh",
			setup: func(m *ctrlauth.MockController) {
				m.EXPECT().Refresh(mock.Anything, "old-refresh").Return(result, nil).Once()
			},
			wantStatus: http.StatusOK,
			wantCookie: true,
		},
		"missing cookie is unauthorized": {
			cookie: "",
			setup: func(m *ctrlauth.MockController) {
				m.EXPECT().Refresh(mock.Anything, "").
					Return(ctrlauth.AuthResult{}, apperror.Unauthorized("INVALID_REFRESH_TOKEN", "invalid refresh token")).Once()
			},
			wantStatus: http.StatusUnauthorized,
		},
	}

	for name, tc := range tcs {
		t.Run(name, func(t *testing.T) {
			ctrl := ctrlauth.NewMockController(t)
			tc.setup(ctrl)

			h := NewAuth(ctrl, false, time.Hour)
			wrapped := httpserver.HandlerErr(h.Refresh)

			rec := httptest.NewRecorder()
			req := httptest.NewRequest(http.MethodPost, "/api/v1/auth/refresh", nil)
			if tc.cookie != "" {
				req.AddCookie(&http.Cookie{Name: refreshCookieName, Value: tc.cookie})
			}
			wrapped.ServeHTTP(rec, req)

			assert.Equal(t, tc.wantStatus, rec.Code)
			if tc.wantCookie {
				cookie := findCookie(rec.Result().Cookies(), refreshCookieName)
				require.NotNil(t, cookie)
				assert.Equal(t, "new-refresh", cookie.Value)
			}
		})
	}
}
