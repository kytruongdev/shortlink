package rest

import (
	"net/http"

	"github.com/kytruongdev/shortlink/internal/infra/httpserver"
)

// Refresh handles POST /api/v1/auth/refresh: rotate the refresh cookie and issue a new access token.
// @Summary      Refresh the access token
// @Description  Read the refresh cookie, rotate it, and return a new access token.
// @Tags         auth
// @Produce      json
// @Success      200  {object}  authResponse
// @Failure      401  {object}  map[string]string
// @Router       /auth/refresh  [post]
func (h *AuthHandler) Refresh(w http.ResponseWriter, r *http.Request) error {
	res, err := h.ctrl.Refresh(r.Context(), refreshTokenFromCookie(r))
	if err != nil {
		return err
	}

	h.setRefreshCookie(w, res.RefreshToken)
	return httpserver.WriteJSON(w, http.StatusOK, toAuthResponse(res, false))
}
