package httpserver_test

import (
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"

	"github.com/kytruongdev/shortlink/internal/infra/httpserver"
)

func TestDecodeJSON(t *testing.T) {
	tcs := map[string]struct {
		body    string
		wantURL string
		wantErr bool
	}{
		"valid":          {body: `{"url":"https://x.com"}`, wantURL: "https://x.com"},
		"unknown field":  {body: `{"url":"x","extra":1}`, wantErr: true},
		"malformed json": {body: `{`, wantErr: true},
		"oversized body": {body: `{"url":"` + strings.Repeat("a", 1<<20+10) + `"}`, wantErr: true},
	}

	for name, tc := range tcs {
		t.Run(name, func(t *testing.T) {
			r := httptest.NewRequest(http.MethodPost, "/", strings.NewReader(tc.body))
			var p struct {
				URL string `json:"url"`
			}

			err := httpserver.DecodeJSON(httptest.NewRecorder(), r, &p)
			if tc.wantErr {
				require.Error(t, err)
				return
			}
			require.NoError(t, err)
			assert.Equal(t, tc.wantURL, p.URL)
		})
	}
}
