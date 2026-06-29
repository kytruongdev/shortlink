package rest

import (
	"bytes"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
	"github.com/stretchr/testify/require"

	ctrllink "github.com/kytruongdev/shortlink/internal/controller/link"
	"github.com/kytruongdev/shortlink/internal/infra/httpserver"
	"github.com/kytruongdev/shortlink/internal/model"
	"github.com/kytruongdev/shortlink/internal/pkg/apperror"
)

func TestDecodeHandler(t *testing.T) {
	const (
		baseURL = "http://short.test"
		code    = "abc1234"
		longURL = "https://example.com/x"
	)

	tcs := map[string]struct {
		body       string
		setup      func(*ctrllink.MockController)
		wantStatus int
		wantBody   map[string]string
	}{
		"decodes full short url": {
			body: `{"short_url":"http://short.test/abc1234"}`,
			setup: func(m *ctrllink.MockController) {
				m.EXPECT().Decode(mock.Anything, code).
					Return(model.Link{OriginalURL: longURL}, nil).Once()
			},
			wantStatus: http.StatusOK,
			wantBody:   map[string]string{"long_url": longURL},
		},
		"decodes bare code": {
			body: `{"short_url":"abc1234"}`,
			setup: func(m *ctrllink.MockController) {
				m.EXPECT().Decode(mock.Anything, code).
					Return(model.Link{OriginalURL: longURL}, nil).Once()
			},
			wantStatus: http.StatusOK,
			wantBody:   map[string]string{"long_url": longURL},
		},
		"not found maps to 404": {
			body: `{"short_url":"http://short.test/missing"}`,
			setup: func(m *ctrllink.MockController) {
				m.EXPECT().Decode(mock.Anything, "missing").
					Return(model.Link{}, apperror.NotFound("NOT_FOUND", "short url not found")).Once()
			},
			wantStatus: http.StatusNotFound,
		},
		"empty short url rejected": {
			body:       `{"short_url":"   "}`,
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

			h := New(ctrl, baseURL)
			wrapped := httpserver.HandlerErr(h.Decode)

			rec := httptest.NewRecorder()
			req := httptest.NewRequest(http.MethodPost, "/api/v1/decode", bytes.NewReader([]byte(tc.body)))
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
