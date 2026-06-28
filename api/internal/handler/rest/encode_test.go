package rest

import (
	"bytes"
	"encoding/json"
	"errors"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
	"github.com/stretchr/testify/require"

	ctrllink "github.com/kytruongdev/shortlink/internal/controller/link"
	"github.com/kytruongdev/shortlink/internal/infra/httpserver"
	"github.com/kytruongdev/shortlink/internal/model"
)

func TestEncodeHandler(t *testing.T) {
	const (
		baseURL  = "http://short.test"
		wantCode = "abc1234"
	)

	tcs := map[string]struct {
		body       string
		setup      func(*ctrllink.MockController)
		wantStatus int
		wantBody   map[string]string
	}{
		"success": {
			body: `{"long_url":"https://example.com"}`,
			setup: func(m *ctrllink.MockController) {
				m.EXPECT().Encode(mock.Anything, "https://example.com").
					Return(model.Link{Code: wantCode}, nil).Once()
			},
			wantStatus: http.StatusOK,
			wantBody:   map[string]string{"short_url": baseURL + "/" + wantCode, "code": wantCode},
		},
		"trims surrounding whitespace before controller": {
			body: `{"long_url":"  https://example.com  "}`,
			setup: func(m *ctrllink.MockController) {
				m.EXPECT().Encode(mock.Anything, "https://example.com").
					Return(model.Link{Code: wantCode}, nil).Once()
			},
			wantStatus: http.StatusOK,
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
		"unknown error maps to 500": {
			body: `{"long_url":"https://example.com"}`,
			setup: func(m *ctrllink.MockController) {
				m.EXPECT().Encode(mock.Anything, mock.Anything).
					Return(model.Link{}, errors.New("boom")).Once()
			},
			wantStatus: http.StatusInternalServerError,
		},
	}

	for name, tc := range tcs {
		t.Run(name, func(t *testing.T) {
			ctrl := ctrllink.NewMockController(t)
			tc.setup(ctrl)

			h := New(ctrl, baseURL)
			wrapped := httpserver.HandlerErr(h.Encode)

			rec := httptest.NewRecorder()
			req := httptest.NewRequest(http.MethodPost, "/api/v1/encode", bytes.NewReader([]byte(tc.body)))
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
