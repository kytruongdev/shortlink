package rest

import (
	"net/http"
	"net/mail"
	"time"

	controllerauth "github.com/kytruongdev/shortlink/internal/controller/auth"
	"github.com/kytruongdev/shortlink/internal/model"
	"github.com/kytruongdev/shortlink/internal/pkg/apperror"
)

const (
	refreshCookieName = "refresh_token"
	refreshCookiePath = "/api/v1/auth"

	codeInvalidBody     = "INVALID_BODY"
	codeInvalidEmail    = "INVALID_EMAIL"
	codeInvalidPassword = "INVALID_PASSWORD"
)

// AuthHandler serves the authentication endpoints under /api/v1/auth.
type AuthHandler struct {
	ctrl         controllerauth.Controller
	cookieSecure bool
	refreshTTL   time.Duration
}

// NewAuth returns an AuthHandler; cookieSecure and refreshTTL shape the refresh cookie.
func NewAuth(ctrl controllerauth.Controller, cookieSecure bool, refreshTTL time.Duration) *AuthHandler {
	return &AuthHandler{ctrl: ctrl, cookieSecure: cookieSecure, refreshTTL: refreshTTL}
}

type authRequest struct {
	Email    string `json:"email"`
	Password string `json:"password"`
}

type authResponse struct {
	AccessToken string        `json:"access_token"`
	ExpiresIn   int           `json:"expires_in"`
	User        *userResponse `json:"user,omitempty"`
}

type userResponse struct {
	ID    string `json:"id"`
	Email string `json:"email"`
}

// toAuthResponse builds the JSON body; user is included only for register/login.
func toAuthResponse(res controllerauth.AuthResult, withUser bool) authResponse {
	out := authResponse{AccessToken: res.AccessToken, ExpiresIn: int(res.ExpiresIn.Seconds())}
	if withUser {
		out.User = &userResponse{ID: res.User.ID.String(), Email: res.User.Email}
	}
	return out
}

// validateCredentials checks the form of sign-up input (handler-level validation).
// It requires a bare email address: ParseAddress alone would also accept display-name
// forms like "Foo <foo@example.com>", whose parsed Address differs from the input.
func validateCredentials(email, password string) error {
	if addr, err := mail.ParseAddress(email); err != nil || addr.Address != email {
		return apperror.BadRequest(codeInvalidEmail, "invalid email address")
	}
	if len(password) < model.MinPasswordLen || len(password) > model.MaxPasswordLen {
		return apperror.BadRequest(codeInvalidPassword, "password must be 8-24 characters")
	}
	return nil
}

// setRefreshCookie writes the refresh token as an httpOnly, same-site cookie.
func (h *AuthHandler) setRefreshCookie(w http.ResponseWriter, token string) {
	http.SetCookie(w, &http.Cookie{
		Name:     refreshCookieName,
		Value:    token,
		Path:     refreshCookiePath,
		HttpOnly: true,
		Secure:   h.cookieSecure,
		SameSite: http.SameSiteStrictMode,
		MaxAge:   int(h.refreshTTL.Seconds()),
	})
}

// clearRefreshCookie expires the refresh cookie.
func (h *AuthHandler) clearRefreshCookie(w http.ResponseWriter) {
	http.SetCookie(w, &http.Cookie{
		Name:     refreshCookieName,
		Value:    "",
		Path:     refreshCookiePath,
		HttpOnly: true,
		Secure:   h.cookieSecure,
		SameSite: http.SameSiteStrictMode,
		MaxAge:   -1,
	})
}

// refreshTokenFromCookie returns the refresh token cookie value, or "" when absent.
func refreshTokenFromCookie(r *http.Request) string {
	c, err := r.Cookie(refreshCookieName)
	if err != nil {
		return ""
	}
	return c.Value
}
