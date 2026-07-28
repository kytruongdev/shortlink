package handler

import (
	"net/http"
	"strings"

	"github.com/google/uuid"

	"github.com/kytruongdev/shortlink/internal/infra/httpserver"
	"github.com/kytruongdev/shortlink/internal/pkg/apperror"
	pkgauth "github.com/kytruongdev/shortlink/internal/pkg/auth"
	"github.com/kytruongdev/shortlink/internal/pkg/authctx"
)

// optionalAuth attaches the authenticated user id to the request context when a
// valid bearer access token is present, and lets anonymous requests through.
func optionalAuth(secret []byte) func(http.Handler) http.Handler {
	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			if token := bearerToken(r); token != "" {
				if claims, err := pkgauth.ParseAccessToken(secret, token); err == nil {
					if id, err := uuid.Parse(claims.UserID); err == nil {
						r = r.WithContext(authctx.WithUser(r.Context(), id))
					}
				}
			}
			next.ServeHTTP(w, r)
		})
	}
}

// requireAuth wraps a handler so it runs only for authenticated requests, else 401.
func requireAuth(fn httpserver.HandlerFunc) httpserver.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) error {
		if _, ok := authctx.UserFromContext(r.Context()); !ok {
			return apperror.Unauthorized("UNAUTHORIZED", "authentication required")
		}
		return fn(w, r)
	}
}

// bearerToken extracts the token from an "Authorization: Bearer <token>" header.
func bearerToken(r *http.Request) string {
	const prefix = "Bearer "
	h := r.Header.Get("Authorization")
	if len(h) < len(prefix) || !strings.EqualFold(h[:len(prefix)], prefix) {
		return ""
	}
	return h[len(prefix):]
}
