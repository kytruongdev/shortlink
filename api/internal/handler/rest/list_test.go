package rest

import (
	"context"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"
	"time"

	"github.com/google/uuid"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
	"github.com/stretchr/testify/require"

	ctrllink "github.com/kytruongdev/shortlink/internal/controller/link"
	"github.com/kytruongdev/shortlink/internal/infra/httpserver"
	"github.com/kytruongdev/shortlink/internal/model"
	"github.com/kytruongdev/shortlink/internal/pkg/authctx"
)

func TestListLinksHandler(t *testing.T) {
	const baseURL = "http://short.test"
	userID := uuid.New()
	created := time.Date(2024, 1, 2, 3, 4, 5, 0, time.UTC)

	tcs := map[string]struct {
		links     []model.Link
		wantCodes []string
	}{
		"returns the user's links": {
			links: []model.Link{
				{Code: "code1", OriginalURL: "https://a.com", CreatedAt: created, ClickCount: 7},
				{Code: "code2", OriginalURL: "https://b.com", CreatedAt: created},
			},
			wantCodes: []string{"code1", "code2"},
		},
		"returns an empty list": {
			links:     []model.Link{},
			wantCodes: []string{},
		},
	}

	for name, tc := range tcs {
		t.Run(name, func(t *testing.T) {
			ctrl := ctrllink.NewMockController(t)
			ctrl.EXPECT().ListByUser(mock.Anything, userID).Return(tc.links, nil).Once()

			h := New(ctrl, baseURL)
			req := httptest.NewRequest(http.MethodGet, "/api/v1/links", nil).
				WithContext(authctx.WithUser(context.Background(), userID))
			rec := httptest.NewRecorder()
			httpserver.HandlerErr(h.ListLinks).ServeHTTP(rec, req)

			assert.Equal(t, http.StatusOK, rec.Code)

			var body listLinksResponse
			require.NoError(t, json.Unmarshal(rec.Body.Bytes(), &body))

			codes := make([]string, len(body.Links))
			for i, l := range body.Links {
				codes[i] = l.Code
			}
			assert.Equal(t, tc.wantCodes, codes)
			if len(tc.links) > 0 {
				assert.Equal(t, baseURL+"/code1", body.Links[0].ShortURL)
				assert.Equal(t, "https://a.com", body.Links[0].LongURL)
				assert.Equal(t, 7, body.Links[0].ClickCount)
			}
		})
	}
}
