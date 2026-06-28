package httpserver

import (
	"context"
	"errors"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/stretchr/testify/assert"
)

type fakePinger struct{ err error }

func (f fakePinger) Ping(context.Context) error { return f.err }

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
			mux := New(fakePinger{tc.pingErr}, nil)

			rec := httptest.NewRecorder()
			mux.ServeHTTP(rec, httptest.NewRequest(http.MethodGet, tc.path, nil))

			assert.Equal(t, tc.wantStatus, rec.Code)
		})
	}
}
