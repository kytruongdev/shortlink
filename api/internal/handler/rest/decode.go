package rest

import (
	"fmt"
	"log/slog"
	"net/http"
	"net/url"
	"path"
	"strings"

	"github.com/kytruongdev/shortlink/internal/infra/httpserver"
	"github.com/kytruongdev/shortlink/internal/pkg/apperror"
)

type decodeRequest struct {
	ShortURL string `json:"short_url"`
}

type decodeResponse struct {
	LongURL string `json:"long_url"`
}

// Decode handles POST /api/v1/decode: short URL (or bare code) in, original URL out.
func (h *Handler) Decode(w http.ResponseWriter, r *http.Request) error {
	slog.Debug("starting Decode")

	const (
		codeInvalidBody     = "INVALID_BODY"
		codeInvalidShortURL = "INVALID_SHORT_URL"
	)

	var req decodeRequest
	if err := httpserver.DecodeJSON(w, r, &req); err != nil {
		slog.Warn("failed to decode json body", slog.String("error", fmt.Sprintf("%+v", err)))
		return apperror.BadRequest(codeInvalidBody, "invalid request body")
	}

	code := codeFromShortURL(req.ShortURL)
	if code == "" {
		return apperror.BadRequest(codeInvalidShortURL, "short url is required")
	}

	link, err := h.ctrl.Decode(r.Context(), code)
	if err != nil {
		return err
	}

	slog.Debug("finished Decode")

	return httpserver.WriteJSON(w, http.StatusOK, decodeResponse{LongURL: link.OriginalURL})
}

// codeFromShortURL extracts the short code from a full short URL or a bare code.
func codeFromShortURL(raw string) string {
	raw = strings.TrimSpace(raw)
	u, err := url.Parse(raw)
	if err != nil || u.Path == "" {
		return raw
	}
	return path.Base(strings.TrimRight(u.Path, "/"))
}
