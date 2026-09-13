package rest

import (
	"bytes"
	"encoding/json"
	"errors"
	"net/http"
	"net/http/httptest"
	"testing"

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

func TestEncodeHandler(t *testing.T) {
	const (
		baseURL  = "http://short.test"
		wantCode = "abc1234"
	)
	userID := uuid.New()

	tcs := map[string]struct {
		body       string
		userID     *uuid.UUID // when set, injected into the request context
		setup      func(*ctrllink.MockController)
		wantStatus int
		wantBody   map[string]string
	}{
		"success passes anonymous caller through": {
			body: `{"long_url":"https://example.com"}`,
			setup: func(m *ctrllink.MockController) {
				m.EXPECT().Encode(mock.Anything, "https://example.com", (*uuid.UUID)(nil), mock.Anything).
					Return(model.Link{Code: wantCode}, nil).Once()
			},
			wantStatus: http.StatusOK,
			wantBody:   map[string]string{"short_url": baseURL + "/" + wantCode, "code": wantCode},
		},
		"authenticated user id is forwarded to the controller": {
			body:   `{"long_url":"https://example.com"}`,
			userID: &userID,
			setup: func(m *ctrllink.MockController) {
				m.EXPECT().Encode(mock.Anything, "https://example.com",
					mock.MatchedBy(func(id *uuid.UUID) bool { return id != nil && *id == userID }), mock.Anything).
					Return(model.Link{Code: wantCode}, nil).Once()
			},
			wantStatus: http.StatusOK,
		},
		"trims surrounding whitespace before controller": {
			body: `{"long_url":"  https://example.com  "}`,
			setup: func(m *ctrllink.MockController) {
				m.EXPECT().Encode(mock.Anything, "https://example.com", mock.Anything, mock.Anything).
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
		"quota exceeded maps to 429": {
			body: `{"long_url":"https://example.com"}`,
			setup: func(m *ctrllink.MockController) {
				m.EXPECT().Encode(mock.Anything, "https://example.com", mock.Anything, mock.Anything).
					Return(model.Link{}, apperror.TooManyRequests("QUOTA_EXCEEDED", "daily link limit reached")).Once()
			},
			wantStatus: http.StatusTooManyRequests,
		},
		"unknown error maps to 500": {
			body: `{"long_url":"https://example.com"}`,
			setup: func(m *ctrllink.MockController) {
				m.EXPECT().Encode(mock.Anything, mock.Anything, mock.Anything, mock.Anything).
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
			if tc.userID != nil {
				req = req.WithContext(authctx.WithUser(req.Context(), *tc.userID))
			}
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
