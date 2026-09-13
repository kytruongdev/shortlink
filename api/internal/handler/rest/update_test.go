package rest

import (
	"bytes"
	"context"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/go-chi/chi/v5"
	"github.com/google/uuid"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
	"github.com/stretchr/testify/require"

	ctrllink "github.com/kytruongdev/shortlink/internal/controller/link"
	"github.com/kytruongdev/shortlink/internal/infra/httpserver"
	"github.com/kytruongdev/shortlink/internal/model"
	"github.com/kytruongdev/shortlink/internal/pkg/apperror"
	"github.com/kytruongdev/shortlink/internal/pkg/authctx"
)

func TestUpdateLinkHandler(t *testing.T) {
	const (
		baseURL = "http://short.test"
		code    = "abc1234"
		newURL  = "https://example.com"
	)
	userID := uuid.New()

	tcs := map[string]struct {
		body       string
		setup      func(*ctrllink.MockController)
		wantStatus int
		wantLong   string
	}{
		"updates returns the link": {
			body: `{"long_url":"https://example.com"}`,
			setup: func(m *ctrllink.MockController) {
				m.EXPECT().UpdateURL(mock.Anything, code, userID, newURL).
					Return(model.Link{Code: code, OriginalURL: newURL}, nil).Once()
			},
			wantStatus: http.StatusOK,
			wantLong:   newURL,
		},
		"conflict maps to 409": {
			body: `{"long_url":"https://example.com"}`,
			setup: func(m *ctrllink.MockController) {
				m.EXPECT().UpdateURL(mock.Anything, code, userID, newURL).
					Return(model.Link{}, apperror.Conflict("URL_TAKEN", "this url is already shortened")).Once()
			},
			wantStatus: http.StatusConflict,
		},
		"invalid url rejected before controller": {
			body:       `{"long_url":"ftp://x"}`,
			setup:      func(*ctrllink.MockController) {},
			wantStatus: http.StatusBadRequest,
		},
		"malformed body maps to 400": {
			body:       `{`,
			setup:      func(*ctrllink.MockController) {},
			wantStatus: http.StatusBadRequest,
		},
	}

	for name, tc := range tcs {
		t.Run(name, func(t *testing.T) {
			ctrl := ctrllink.NewMockController(t)
			tc.setup(ctrl)

			wrapped := httpserver.HandlerErr(New(ctrl, baseURL).UpdateLink)

			rctx := chi.NewRouteContext()
			rctx.URLParams.Add("code", code)
			ctx := authctx.WithUser(context.WithValue(context.Background(), chi.RouteCtxKey, rctx), userID)
			req := httptest.NewRequest(http.MethodPatch, "/api/v1/links/"+code, bytes.NewReader([]byte(tc.body))).WithContext(ctx)

			rec := httptest.NewRecorder()
			wrapped.ServeHTTP(rec, req)

			assert.Equal(t, tc.wantStatus, rec.Code)
			if tc.wantLong != "" {
				var got map[string]any
				require.NoError(t, json.Unmarshal(rec.Body.Bytes(), &got))
				assert.Equal(t, tc.wantLong, got["long_url"])
			}
		})
	}
}
