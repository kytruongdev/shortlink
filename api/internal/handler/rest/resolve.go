package rest

import (
	"net/http"

	"github.com/go-chi/chi/v5"

	"github.com/kytruongdev/shortlink/internal/infra/httpserver"
	"github.com/kytruongdev/shortlink/internal/pkg/apperror"
)

type resolveResponse struct {
	LongURL string `json:"long_url"`
}

// Resolve handles GET /api/v1/{code}: resolve a short code to its original URL.
// @Summary      Resolve a short code
// @Description  Resolve a short code to the original URL; the frontend uses this to redirect the visitor.
// @Tags         links
// @Produce      json
// @Param        code  path      string  true  "Short code"
// @Success      200   {object}  resolveResponse
// @Failure      400   {object}  map[string]string
// @Failure      404   {object}  map[string]string
// @Router       /{code}  [get]
func (h *Handler) Resolve(w http.ResponseWriter, r *http.Request) error {
	const codeInvalidShortURL = "INVALID_SHORT_URL"

	code := chi.URLParam(r, "code")
	if code == "" {
		return apperror.BadRequest(codeInvalidShortURL, "short code is required")
	}

	link, err := h.ctrl.Decode(r.Context(), code)
	if err != nil {
		return err
	}

	return httpserver.WriteJSON(w, http.StatusOK, resolveResponse{LongURL: link.OriginalURL})
}
