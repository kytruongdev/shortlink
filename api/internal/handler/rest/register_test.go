package rest

import (
	"bytes"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"strings"
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

// findCookie returns the named cookie from a response, or nil.
func findCookie(cookies []*http.Cookie, name string) *http.Cookie {
	for _, c := range cookies {
		if c.Name == name {
			return c
		}
	}
	return nil
}

func TestRegisterHandler(t *testing.T) {
	result := ctrlauth.AuthResult{
		User:         model.User{ID: uuid.New(), Email: "new@example.com"},
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
		"registers, returns token and sets refresh cookie": {
			body: `{"email":"new@example.com","password":"password123"}`,
			setup: func(m *ctrlauth.MockController) {
				m.EXPECT().Register(mock.Anything, "new@example.com", "password123").Return(result, nil).Once()
			},
			wantStatus: http.StatusOK,
			wantCookie: true,
		},
		"invalid email rejected before controller": {
			body:       `{"email":"not-an-email","password":"password123"}`,
			setup:      func(*ctrlauth.MockController) {},
			wantStatus: http.StatusBadRequest,
		},
		"display-name email rejected": {
			body:       `{"email":"Foo <foo@example.com>","password":"password123"}`,
			setup:      func(*ctrlauth.MockController) {},
			wantStatus: http.StatusBadRequest,
		},
		"short password rejected": {
			body:       `{"email":"ok@example.com","password":"short"}`,
			setup:      func(*ctrlauth.MockController) {},
			wantStatus: http.StatusBadRequest,
		},
		"too long password rejected": {
			body:       `{"email":"ok@example.com","password":"` + strings.Repeat("a", model.MaxPasswordLen+1) + `"}`,
			setup:      func(*ctrlauth.MockController) {},
			wantStatus: http.StatusBadRequest,
		},
		"duplicate email surfaces as conflict": {
			body: `{"email":"dupe@example.com","password":"password123"}`,
			setup: func(m *ctrlauth.MockController) {
				m.EXPECT().Register(mock.Anything, "dupe@example.com", "password123").
					Return(ctrlauth.AuthResult{}, apperror.Conflict("EMAIL_TAKEN", "email already registered")).Once()
			},
			wantStatus: http.StatusConflict,
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
			wrapped := httpserver.HandlerErr(h.Register)

			rec := httptest.NewRecorder()
			req := httptest.NewRequest(http.MethodPost, "/api/v1/auth/register", bytes.NewReader([]byte(tc.body)))
			wrapped.ServeHTTP(rec, req)

			assert.Equal(t, tc.wantStatus, rec.Code)
			if !tc.wantCookie {
				return
			}

			var body map[string]any
			require.NoError(t, json.Unmarshal(rec.Body.Bytes(), &body))
			assert.Equal(t, "access-token", body["access_token"])

			cookie := findCookie(rec.Result().Cookies(), refreshCookieName)
			require.NotNil(t, cookie)
			assert.Equal(t, "refresh-raw", cookie.Value)
			assert.True(t, cookie.HttpOnly)
			assert.Equal(t, http.SameSiteStrictMode, cookie.SameSite)
			assert.Equal(t, refreshCookiePath, cookie.Path)
		})
	}
}
