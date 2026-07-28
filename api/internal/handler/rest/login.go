package rest

import (
	"fmt"
	"log/slog"
	"net/http"
	"strings"

	"github.com/kytruongdev/shortlink/internal/infra/httpserver"
	"github.com/kytruongdev/shortlink/internal/pkg/apperror"
)

// Login handles POST /api/v1/auth/login: verify credentials and issue tokens.
// @Summary      Log in
// @Description  Verify credentials, returning an access token and setting an httpOnly refresh cookie.
// @Tags         auth
// @Accept       json
// @Produce      json
// @Param        request  body      authRequest  true  "Email and password"
// @Success      200      {object}  authResponse
// @Failure      400      {object}  map[string]string
// @Failure      401      {object}  map[string]string
// @Router       /auth/login  [post]
func (h *AuthHandler) Login(w http.ResponseWriter, r *http.Request) error {
	var req authRequest
	if err := httpserver.DecodeJSON(w, r, &req); err != nil {
		slog.Warn("failed to decode json body", slog.String("error", fmt.Sprintf("%+v", err)))
		return apperror.BadRequest(codeInvalidBody, "invalid request body")
	}

	res, err := h.ctrl.Login(r.Context(), strings.TrimSpace(req.Email), req.Password)
	if err != nil {
		return err
	}

	h.setRefreshCookie(w, res.RefreshToken)
	return httpserver.WriteJSON(w, http.StatusOK, toAuthResponse(res, true))
}
