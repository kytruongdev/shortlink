package rest

import (
	"context"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/go-chi/chi/v5"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
	"github.com/stretchr/testify/require"

	ctrllink "github.com/kytruongdev/shortlink/internal/controller/link"
	"github.com/kytruongdev/shortlink/internal/infra/httpserver"
	"github.com/kytruongdev/shortlink/internal/model"
	"github.com/kytruongdev/shortlink/internal/pkg/apperror"
)

func TestResolveHandler(t *testing.T) {
	const (
		baseURL = "http://short.test"
		code    = "abc1234"
		longURL = "https://example.com/x"
	)

	tcs := map[string]struct {
		code       string
		setup      func(*ctrllink.MockController)
		wantStatus int
		wantBody   map[string]string
	}{
		"found returns long url": {
			code: code,
			setup: func(m *ctrllink.MockController) {
				m.EXPECT().Decode(mock.Anything, code).
					Return(model.Link{OriginalURL: longURL}, nil).Once()
			},
			wantStatus: http.StatusOK,
			wantBody:   map[string]string{"long_url": longURL},
		},
		"not found maps to 404": {
			code: "missing",
			setup: func(m *ctrllink.MockController) {
				m.EXPECT().Decode(mock.Anything, "missing").
					Return(model.Link{}, apperror.NotFound("NOT_FOUND", "short url not found")).Once()
			},
			wantStatus: http.StatusNotFound,
		},
		"empty code rejected": {
			code:       "",
			setup:      func(*ctrllink.MockController) {},
			wantStatus: http.StatusBadRequest,
		},
	}

	for name, tc := range tcs {
		t.Run(name, func(t *testing.T) {
			ctrl := ctrllink.NewMockController(t)
			tc.setup(ctrl)

			h := New(ctrl, baseURL)
			wrapped := httpserver.HandlerErr(h.Resolve)

			rctx := chi.NewRouteContext()
			rctx.URLParams.Add("code", tc.code)

			rec := httptest.NewRecorder()
			req := httptest.NewRequest(http.MethodGet, "/api/v1/"+tc.code, nil)
			req = req.WithContext(context.WithValue(req.Context(), chi.RouteCtxKey, rctx))
			wrapped.ServeHTTP(rec, req)

			assert.Equal(t, tc.wantStatus, rec.Code)
			if tc.wantBody != nil {
				var got map[string]string
				require.NoError(t, json.Unmarshal(rec.Body.Bytes(), &got))
				for k, v := range tc.wantBody {
					assert.Equal(t, v, got[k])
				}
			}
		})
	}
}
