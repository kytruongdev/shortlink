package httpserver_test

import (
	"context"
	"errors"
	"io"
	"log/slog"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/stretchr/testify/assert"

	"github.com/kytruongdev/shortlink/internal/infra/httpserver"
)

type fakePinger struct{ err error }

func (f fakePinger) Ping(context.Context) error { return f.err }

func discardLogger() *slog.Logger {
	return slog.New(slog.NewTextHandler(io.Discard, nil))
}

func TestHealthRoutes(t *testing.T) {
	tcs := map[string]struct {
		pingErr    error
		path       string
		wantStatus int
	}{
		"healthz ok":                   {pingErr: nil, path: "/healthz", wantStatus: http.StatusOK},
		"readyz ok when db up":         {pingErr: nil, path: "/readyz", wantStatus: http.StatusOK},
		"healthz ok even when db down": {pingErr: errors.New("db down"), path: "/healthz", wantStatus: http.StatusOK},
		"readyz 503 when db down":      {pingErr: errors.New("db down"), path: "/readyz", wantStatus: http.StatusServiceUnavailable},
	}

	for name, tc := range tcs {
		t.Run(name, func(t *testing.T) {
			mux := httpserver.New(discardLogger(), fakePinger{tc.pingErr}, nil)

			rec := httptest.NewRecorder()
			mux.ServeHTTP(rec, httptest.NewRequest(http.MethodGet, tc.path, nil))

			assert.Equal(t, tc.wantStatus, rec.Code)
		})
	}
}
