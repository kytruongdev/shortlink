package handler

import (
	"net/http"
	"net/http/httptest"
	"testing"
	"time"

	"github.com/google/uuid"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"

	pkgauth "github.com/kytruongdev/shortlink/internal/pkg/auth"
	"github.com/kytruongdev/shortlink/internal/pkg/authctx"
)

func TestOptionalAuth(t *testing.T) {
	secret := []byte("test-secret-test-secret-test-secret")
	userID := uuid.New()
	validToken, err := pkgauth.SignAccessToken(secret, userID.String(), time.Minute)
	require.NoError(t, err)

	tcs := map[string]struct {
		authHeader string
		wantUser   bool
	}{
		"valid bearer token sets user":   {authHeader: "Bearer " + validToken, wantUser: true},
		"missing header stays anonymous": {authHeader: "", wantUser: false},
		"invalid token stays anonymous":  {authHeader: "Bearer not-a-token", wantUser: false},
		"non-bearer header ignored":      {authHeader: validToken, wantUser: false},
	}

	for name, tc := range tcs {
		t.Run(name, func(t *testing.T) {
			var gotID uuid.UUID
			var gotUser bool
			next := http.HandlerFunc(func(_ http.ResponseWriter, r *http.Request) {
				gotID, gotUser = authctx.UserFromContext(r.Context())
			})

			req := httptest.NewRequest(http.MethodGet, "/api/v1/encode", nil)
			if tc.authHeader != "" {
				req.Header.Set("Authorization", tc.authHeader)
			}
			optionalAuth(secret)(next).ServeHTTP(httptest.NewRecorder(), req)

			assert.Equal(t, tc.wantUser, gotUser)
			if tc.wantUser {
				assert.Equal(t, userID, gotID)
			}
		})
	}
}

func TestRequireAuth(t *testing.T) {
	tcs := map[string]struct {
		authenticated bool
		wantStatus    int
		wantCalled    bool
	}{
		"blocks anonymous with 401": {authenticated: false, wantStatus: http.StatusUnauthorized, wantCalled: false},
		"allows authenticated":      {authenticated: true, wantStatus: http.StatusOK, wantCalled: true},
	}

	for name, tc := range tcs {
		t.Run(name, func(t *testing.T) {
			called := false
			next := http.HandlerFunc(func(_ http.ResponseWriter, _ *http.Request) {
				called = true
			})

			req := httptest.NewRequest(http.MethodGet, "/api/v1/links", nil)
			if tc.authenticated {
				req = req.WithContext(authctx.WithUser(req.Context(), uuid.New()))
			}
			rec := httptest.NewRecorder()
			requireAuth(next).ServeHTTP(rec, req)

			assert.Equal(t, tc.wantStatus, rec.Code)
			assert.Equal(t, tc.wantCalled, called)
		})
	}
}
