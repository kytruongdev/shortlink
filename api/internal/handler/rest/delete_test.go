package rest

import (
	"context"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/go-chi/chi/v5"
	"github.com/google/uuid"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"

	ctrllink "github.com/kytruongdev/shortlink/internal/controller/link"
	"github.com/kytruongdev/shortlink/internal/infra/httpserver"
	"github.com/kytruongdev/shortlink/internal/pkg/apperror"
	"github.com/kytruongdev/shortlink/internal/pkg/authctx"
)

func TestDeleteLinkHandler(t *testing.T) {
	const (
		baseURL = "http://short.test"
		code    = "abc1234"
	)
	userID := uuid.New()

	tcs := map[string]struct {
		setup      func(*ctrllink.MockController)
		wantStatus int
	}{
		"deletes returns 204": {
			setup: func(m *ctrllink.MockController) {
				m.EXPECT().Delete(mock.Anything, code, userID).Return(nil).Once()
			},
			wantStatus: http.StatusNoContent,
		},
		"not found maps to 404": {
			setup: func(m *ctrllink.MockController) {
				m.EXPECT().Delete(mock.Anything, code, userID).
					Return(apperror.NotFound("NOT_FOUND", "link not found")).Once()
			},
			wantStatus: http.StatusNotFound,
		},
	}

	for name, tc := range tcs {
		t.Run(name, func(t *testing.T) {
			ctrl := ctrllink.NewMockController(t)
			tc.setup(ctrl)

			wrapped := httpserver.HandlerErr(New(ctrl, baseURL).DeleteLink)

			rctx := chi.NewRouteContext()
			rctx.URLParams.Add("code", code)
			ctx := authctx.WithUser(context.WithValue(context.Background(), chi.RouteCtxKey, rctx), userID)
			req := httptest.NewRequest(http.MethodDelete, "/api/v1/links/"+code, nil).WithContext(ctx)

			rec := httptest.NewRecorder()
			wrapped.ServeHTTP(rec, req)

			assert.Equal(t, tc.wantStatus, rec.Code)
		})
	}
}
