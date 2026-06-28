package urlshortener

import (
	"strings"
	"testing"

	"github.com/stretchr/testify/require"
)

func TestValidate(t *testing.T) {
	const maxLen = 2048

	tcs := map[string]struct {
		url     string
		wantErr bool
	}{
		"valid http":          {url: "http://example.com", wantErr: false},
		"valid https path":    {url: "https://example.com/a/b?x=1", wantErr: false},
		"valid subdomain":     {url: "https://www.example.co.uk/x", wantErr: false},
		"valid trailing dot":  {url: "https://example.com./x", wantErr: false},
		"empty":               {url: "", wantErr: true},
		"whitespace only":     {url: "   ", wantErr: true},
		"too long":            {url: "https://e.com/" + strings.Repeat("a", maxLen), wantErr: true},
		"ftp scheme":          {url: "ftp://example.com", wantErr: true},
		"javascript scheme":   {url: "javascript:alert(1)", wantErr: true},
		"no host":             {url: "http://", wantErr: true},
		"space in host":       {url: "http://exa mple.com", wantErr: true},
		"userinfo phishing":   {url: "https://www.youtube@com.vn/x", wantErr: true},
		"bang in host":        {url: "https://youtube.com.vn!/x", wantErr: true},
		"single label no dot": {url: "http://com", wantErr: true},
		"localhost":           {url: "http://localhost:8080", wantErr: true},
		"public ipv4":         {url: "http://8.8.8.8", wantErr: true},
		"private ipv4":        {url: "http://192.168.1.1/x", wantErr: true},
		"ipv6 loopback":       {url: "http://[::1]/x", wantErr: true},
	}

	for name, tc := range tcs {
		t.Run(name, func(t *testing.T) {
			err := Validate(tc.url, maxLen)
			if tc.wantErr {
				require.Error(t, err)
				return
			}
			require.NoError(t, err)
		})
	}
}
