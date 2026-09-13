package rest

import (
	"net/http"

	"github.com/go-chi/chi/v5"

	"github.com/kytruongdev/shortlink/internal/pkg/authctx"
)

// DeleteLink handles DELETE /api/v1/links/{code}: remove the current user's link.
// @Summary   Delete a link
// @Tags      links
// @Security  BearerAuth
// @Param     code  path  string  true  "Short code"
// @Success   204   "No Content"
// @Failure   401   {object}  map[string]string
// @Failure   404   {object}  map[string]string
// @Router    /links/{code}  [delete]
func (h *Handler) DeleteLink(w http.ResponseWriter, r *http.Request) error {
	userID, _ := authctx.UserFromContext(r.Context())

	if err := h.ctrl.Delete(r.Context(), chi.URLParam(r, "code"), userID); err != nil {
		return err
	}

	w.WriteHeader(http.StatusNoContent)
	return nil
}
