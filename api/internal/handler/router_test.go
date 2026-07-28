package handler

import (
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
	"github.com/stretchr/testify/require"

	ctrllink "github.com/kytruongdev/shortlink/internal/controller/link"
	"github.com/kytruongdev/shortlink/internal/handler/rest"
	"github.com/kytruongdev/shortlink/internal/infra/httpserver"
	"github.com/kytruongdev/shortlink/internal/model"
)

// TestRoutes verifies GET /api/v1/{code} is registered and binds the code param
// through the fully wired router.
func TestRoutes(t *testing.T) {
	const (
		baseURL = "http://short.test"
		code    = "abc1234"
		longURL = "https://example.com/x"
	)

	ctrl := ctrllink.NewMockController(t)
	ctrl.EXPECT().Decode(mock.Anything, code).
		Return(model.Link{OriginalURL: longURL}, nil).Once()

	srv := httpserver.New(nil, New(rest.New(ctrl, baseURL)).Routes)

	rec := httptest.NewRecorder()
	req := httptest.NewRequest(http.MethodGet, "/api/v1/"+code, nil)
	srv.ServeHTTP(rec, req)

	assert.Equal(t, http.StatusOK, rec.Code)

	var got map[string]string
	require.NoError(t, json.Unmarshal(rec.Body.Bytes(), &got))
	assert.Equal(t, longURL, got["long_url"])
}
