package httpserver_test

import (
	"errors"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/stretchr/testify/assert"

	"github.com/kytruongdev/shortlink/internal/infra/httpserver"
)

func TestHandlerErr(t *testing.T) {
	tcs := map[string]struct {
		fn               httpserver.HandlerFunc
		wantStatus       int
		wantBodyContains string
	}{
		"returned error becomes 500": {
			fn:               func(http.ResponseWriter, *http.Request) error { return errors.New("boom") },
			wantStatus:       http.StatusInternalServerError,
			wantBodyContains: "error",
		},
		"nil error preserves handler response": {
			fn:         func(w http.ResponseWriter, _ *http.Request) error { w.WriteHeader(http.StatusNoContent); return nil },
			wantStatus: http.StatusNoContent,
		},
	}

	for name, tc := range tcs {
		t.Run(name, func(t *testing.T) {
			h := httpserver.HandlerErr(discardLogger(), tc.fn)

			rec := httptest.NewRecorder()
			h.ServeHTTP(rec, httptest.NewRequest(http.MethodGet, "/", nil))

			assert.Equal(t, tc.wantStatus, rec.Code)
			if tc.wantBodyContains != "" {
				assert.Contains(t, rec.Body.String(), tc.wantBodyContains)
			}
		})
	}
}
