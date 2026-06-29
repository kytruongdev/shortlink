package urlshortener

import (
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestGenerate(t *testing.T) {
	tcs := map[string]struct {
		n int
	}{
		"length 1":  {n: 1},
		"length 7":  {n: 7},
		"length 20": {n: 20},
	}

	for name, tc := range tcs {
		t.Run(name, func(t *testing.T) {
			code, err := Generate(tc.n)
			require.NoError(t, err)
			assert.Len(t, code, tc.n)
			for _, c := range code {
				assert.Contains(t, alphabet, string(c))
			}
		})
	}
}
