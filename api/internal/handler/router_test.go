package handler

import (
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
	"github.com/stretchr/testify/require"

	ctrlauth "github.com/kytruongdev/shortlink/internal/controller/auth"
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
	ctrl.EXPECT().Resolve(mock.Anything, code).
		Return(model.Link{OriginalURL: longURL}, nil).Once()

	authHandler := rest.NewAuth(ctrlauth.NewMockController(t), false, time.Hour)
	srv := httpserver.New(nil, nil, New(rest.New(ctrl, baseURL), authHandler, []byte("test-secret")).Routes)

	rec := httptest.NewRecorder()
	req := httptest.NewRequest(http.MethodGet, "/api/v1/"+code, nil)
	srv.ServeHTTP(rec, req)

	assert.Equal(t, http.StatusOK, rec.Code)

	var got map[string]string
	require.NoError(t, json.Unmarshal(rec.Body.Bytes(), &got))
	assert.Equal(t, longURL, got["long_url"])
}

// TestListLinksRequiresAuth verifies GET /api/v1/links is guarded by requireAuth.
func TestListLinksRequiresAuth(t *testing.T) {
	authHandler := rest.NewAuth(ctrlauth.NewMockController(t), false, time.Hour)
	srv := httpserver.New(nil, nil, New(rest.New(ctrllink.NewMockController(t), "http://short.test"), authHandler, []byte("test-secret")).Routes)

	rec := httptest.NewRecorder()
	req := httptest.NewRequest(http.MethodGet, "/api/v1/links", nil)
	srv.ServeHTTP(rec, req)

	assert.Equal(t, http.StatusUnauthorized, rec.Code)
}
