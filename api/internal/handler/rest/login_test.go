package rest

import (
	"bytes"
	"net/http"
	"net/http/httptest"
	"testing"
	"time"

	"github.com/google/uuid"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
	"github.com/stretchr/testify/require"

	ctrlauth "github.com/kytruongdev/shortlink/internal/controller/auth"
	"github.com/kytruongdev/shortlink/internal/infra/httpserver"
	"github.com/kytruongdev/shortlink/internal/model"
	"github.com/kytruongdev/shortlink/internal/pkg/apperror"
)

func TestLoginHandler(t *testing.T) {
	result := ctrlauth.AuthResult{
		User:         model.User{ID: uuid.New(), Email: "user@example.com"},
		AccessToken:  "access-token",
		RefreshToken: "refresh-raw",
		ExpiresIn:    5 * time.Minute,
	}

	tcs := map[string]struct {
		body       string
		setup      func(*ctrlauth.MockController)
		wantStatus int
		wantCookie bool
	}{
		"logs in and sets refresh cookie": {
			body: `{"email":"user@example.com","password":"password123"}`,
			setup: func(m *ctrlauth.MockController) {
				m.EXPECT().Login(mock.Anything, "user@example.com", "password123").Return(result, nil).Once()
			},
			wantStatus: http.StatusOK,
			wantCookie: true,
		},
		"invalid credentials propagate as 401": {
			body: `{"email":"user@example.com","password":"wrong-password"}`,
			setup: func(m *ctrlauth.MockController) {
				m.EXPECT().Login(mock.Anything, "user@example.com", "wrong-password").
					Return(ctrlauth.AuthResult{}, apperror.Unauthorized("INVALID_CREDENTIALS", "invalid email or password")).Once()
			},
			wantStatus: http.StatusUnauthorized,
		},
		"malformed body rejected": {
			body:       `{`,
			setup:      func(*ctrlauth.MockController) {},
			wantStatus: http.StatusBadRequest,
		},
	}

	for name, tc := range tcs {
		t.Run(name, func(t *testing.T) {
			ctrl := ctrlauth.NewMockController(t)
			tc.setup(ctrl)

			h := NewAuth(ctrl, false, time.Hour)
			wrapped := httpserver.HandlerErr(h.Login)

			rec := httptest.NewRecorder()
			req := httptest.NewRequest(http.MethodPost, "/api/v1/auth/login", bytes.NewReader([]byte(tc.body)))
			wrapped.ServeHTTP(rec, req)

			assert.Equal(t, tc.wantStatus, rec.Code)
			if tc.wantCookie {
				cookie := findCookie(rec.Result().Cookies(), refreshCookieName)
				require.NotNil(t, cookie)
				assert.Equal(t, "refresh-raw", cookie.Value)
				assert.True(t, cookie.HttpOnly)
			}
		})
	}
}
