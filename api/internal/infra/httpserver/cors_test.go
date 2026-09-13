package httpserver

import (
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestCORSPreflight(t *testing.T) {
	const origin = "http://localhost:5173"
	mux := New(nil, []string{origin}, nil)

	tcs := map[string]struct {
		origin     string
		wantOrigin string
	}{
		"allowed origin is echoed with credentials": {origin: origin, wantOrigin: origin},
		"disallowed origin is not echoed":           {origin: "http://evil.test", wantOrigin: ""},
	}

	for name, tc := range tcs {
		t.Run(name, func(t *testing.T) {
			rec := httptest.NewRecorder()
			req := httptest.NewRequest(http.MethodOptions, "/api/v1/encode", nil)
			req.Header.Set("Origin", tc.origin)
			req.Header.Set("Access-Control-Request-Method", http.MethodPost)
			mux.ServeHTTP(rec, req)

			assert.Equal(t, tc.wantOrigin, rec.Header().Get("Access-Control-Allow-Origin"))
			if tc.wantOrigin != "" {
				assert.Equal(t, "true", rec.Header().Get("Access-Control-Allow-Credentials"))
			}
		})
	}
}
