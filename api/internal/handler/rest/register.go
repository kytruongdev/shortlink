package rest

import (
	"fmt"
	"log/slog"
	"net/http"
	"strings"

	"github.com/kytruongdev/shortlink/internal/infra/httpserver"
	"github.com/kytruongdev/shortlink/internal/pkg/apperror"
)

// Register handles POST /api/v1/auth/register: create an account and auto-login.
// @Summary      Register a new account
// @Description  Create an account, returning an access token and setting an httpOnly refresh cookie.
// @Tags         auth
// @Accept       json
// @Produce      json
// @Param        request  body      authRequest  true  "Email and password"
// @Success      200      {object}  authResponse
// @Failure      400      {object}  map[string]string
// @Failure      409      {object}  map[string]string
// @Router       /auth/register  [post]
func (h *AuthHandler) Register(w http.ResponseWriter, r *http.Request) error {
	var req authRequest
	if err := httpserver.DecodeJSON(w, r, &req); err != nil {
		slog.Warn("failed to decode json body", slog.String("error", fmt.Sprintf("%+v", err)))
		return apperror.BadRequest(codeInvalidBody, "invalid request body")
	}

	email := strings.TrimSpace(req.Email)
	if err := validateCredentials(email, req.Password); err != nil {
		return err
	}

	res, err := h.ctrl.Register(r.Context(), email, req.Password)
	if err != nil {
		return err
	}

	h.setRefreshCookie(w, res.RefreshToken)
	return httpserver.WriteJSON(w, http.StatusOK, toAuthResponse(res, true))
}
