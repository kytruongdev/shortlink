package rest

import (
	"fmt"
	"log/slog"
	"net/http"
	"strings"

	"github.com/kytruongdev/shortlink/internal/infra/httpserver"
	"github.com/kytruongdev/shortlink/internal/model"
	"github.com/kytruongdev/shortlink/internal/pkg/apperror"
	"github.com/kytruongdev/shortlink/internal/pkg/urlshortener"
)

type encodeRequest struct {
	LongURL string `json:"long_url"`
}

type encodeResponse struct {
	ShortURL string `json:"short_url"`
	Code     string `json:"code"`
}

// Encode handles POST /api/v1/encode.
// @Summary      Shorten a URL
// @Description  Encode a long URL into a short code and persist the mapping.
// @Tags         links
// @Accept       json
// @Produce      json
// @Param        request  body      encodeRequest  true  "URL to shorten"
// @Success      200      {object}  encodeResponse
// @Failure      400      {object}  map[string]string
// @Failure      500      {object}  map[string]string
// @Router       /encode  [post]
func (h *Handler) Encode(w http.ResponseWriter, r *http.Request) error {
	slog.Debug("starting Encode")

	const (
		codeInvalidBody = "INVALID_BODY"
		codeInvalidURL  = "INVALID_URL"
	)

	var req encodeRequest
	if err := httpserver.DecodeJSON(w, r, &req); err != nil {
		slog.Warn("failed to decode json body", slog.String("error", fmt.Sprintf("%+v", err)))
		return apperror.BadRequest(codeInvalidBody, "invalid request body")
	}

	longURL := strings.TrimSpace(req.LongURL)
	if err := urlshortener.Validate(longURL, model.MaxURLLen); err != nil {
		return apperror.BadRequest(codeInvalidURL, "invalid url")
	}

	link, err := h.ctrl.Encode(r.Context(), longURL)
	if err != nil {
		return err
	}

	slog.Debug("finished Encode")

	return httpserver.WriteJSON(w, http.StatusOK, encodeResponse{
		ShortURL: h.baseURL + "/" + link.Code,
		Code:     link.Code,
	})
}
