package urlshortener

import (
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestNormalize(t *testing.T) {
	tcs := map[string]struct {
		in   string
		want string
	}{
		"lowercase scheme and host": {in: "HTTP://Example.COM/Path", want: "http://example.com/Path"},
		"strip default http port":   {in: "http://example.com:80/x", want: "http://example.com/x"},
		"strip default https port":  {in: "https://example.com:443/x", want: "https://example.com/x"},
		"keep non-default port":     {in: "http://example.com:8080/x", want: "http://example.com:8080/x"},
		"drop empty fragment":       {in: "http://example.com/x#", want: "http://example.com/x"},
		"keep non-empty fragment":   {in: "http://example.com/x#sec", want: "http://example.com/x#sec"},
		"keep path and query order": {in: "http://example.com/a?b=2&a=1", want: "http://example.com/a?b=2&a=1"},
	}

	for name, tc := range tcs {
		t.Run(name, func(t *testing.T) {
			got, err := Normalize(tc.in)
			require.NoError(t, err)
			assert.Equal(t, tc.want, got)
		})
	}
}
