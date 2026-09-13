package rest

import (
	"net/http"
	"time"

	"github.com/kytruongdev/shortlink/internal/infra/httpserver"
	"github.com/kytruongdev/shortlink/internal/model"
	"github.com/kytruongdev/shortlink/internal/pkg/authctx"
)

type linkResponse struct {
	Code       string    `json:"code"`
	ShortURL   string    `json:"short_url"`
	LongURL    string    `json:"long_url"`
	CreatedAt  time.Time `json:"created_at"`
	ClickCount int       `json:"click_count"`
}

type listLinksResponse struct {
	Links []linkResponse `json:"links"`
}

func (h *Handler) toLinkResponse(l model.Link) linkResponse {
	return linkResponse{
		Code:       l.Code,
		ShortURL:   h.baseURL + "/" + l.Code,
		LongURL:    l.OriginalURL,
		CreatedAt:  l.CreatedAt,
		ClickCount: l.ClickCount,
	}
}

// ListLinks handles GET /api/v1/links: the current user's links, newest first.
// @Summary      List my links
// @Description  Return the authenticated user's links, newest first.
// @Tags         links
// @Produce      json
// @Security     BearerAuth
// @Success      200  {object}  listLinksResponse
// @Failure      401  {object}  map[string]string
// @Router       /links  [get]
func (h *Handler) ListLinks(w http.ResponseWriter, r *http.Request) error {
	// requireAuth guarantees a user is present in the context.
	userID, _ := authctx.UserFromContext(r.Context())

	links, err := h.ctrl.ListByUser(r.Context(), userID)
	if err != nil {
		return err
	}

	out := make([]linkResponse, 0, len(links))
	for _, l := range links {
		out = append(out, h.toLinkResponse(l))
	}

	return httpserver.WriteJSON(w, http.StatusOK, listLinksResponse{Links: out})
}
