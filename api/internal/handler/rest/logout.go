package rest

import (
	"net/http"
)

// Logout handles POST /api/v1/auth/logout: revoke the refresh token and clear the cookie.
// @Summary      Log out
// @Description  Revoke the presented refresh token and clear the refresh cookie.
// @Tags         auth
// @Success      204  "No Content"
// @Router       /auth/logout  [post]
func (h *AuthHandler) Logout(w http.ResponseWriter, r *http.Request) error {
	if err := h.ctrl.Logout(r.Context(), refreshTokenFromCookie(r)); err != nil {
		return err
	}

	h.clearRefreshCookie(w)
	w.WriteHeader(http.StatusNoContent)
	return nil
}
