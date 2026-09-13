package rest

import (
	"fmt"
	"log/slog"
	"net/http"
	"strings"

	"github.com/go-chi/chi/v5"

	"github.com/kytruongdev/shortlink/internal/infra/httpserver"
	"github.com/kytruongdev/shortlink/internal/model"
	"github.com/kytruongdev/shortlink/internal/pkg/apperror"
	"github.com/kytruongdev/shortlink/internal/pkg/authctx"
	"github.com/kytruongdev/shortlink/internal/pkg/urlshortener"
)

type updateLinkRequest struct {
	LongURL string `json:"long_url"`
}

// UpdateLink handles PATCH /api/v1/links/{code}: change the current user's link destination.
// @Summary   Edit a link's destination
// @Tags      links
// @Accept    json
// @Produce   json
// @Security  BearerAuth
// @Param     code     path      string             true  "Short code"
// @Param     request  body      updateLinkRequest  true  "New long URL"
// @Success   200      {object}  linkResponse
// @Failure   400      {object}  map[string]string
// @Failure   401      {object}  map[string]string
// @Failure   404      {object}  map[string]string
// @Failure   409      {object}  map[string]string
// @Router    /links/{code}  [patch]
func (h *Handler) UpdateLink(w http.ResponseWriter, r *http.Request) error {
	const codeInvalidURL = "INVALID_URL"

	var req updateLinkRequest
	if err := httpserver.DecodeJSON(w, r, &req); err != nil {
		slog.Warn("failed to decode json body", slog.String("error", fmt.Sprintf("%+v", err)))
		return apperror.BadRequest(codeInvalidBody, "invalid request body")
	}

	longURL := strings.TrimSpace(req.LongURL)
	if err := urlshortener.Validate(longURL, model.MaxURLLen); err != nil {
		return apperror.BadRequest(codeInvalidURL, "invalid url")
	}

	userID, _ := authctx.UserFromContext(r.Context())
	link, err := h.ctrl.UpdateURL(r.Context(), chi.URLParam(r, "code"), userID, longURL)
	if err != nil {
		return err
	}

	return httpserver.WriteJSON(w, http.StatusOK, h.toLinkResponse(link))
}
