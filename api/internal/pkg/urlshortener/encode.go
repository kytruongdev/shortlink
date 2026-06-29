package urlshortener

import (
	"crypto/rand"

	pkgerrors "github.com/pkg/errors"
)

const alphabet = "0123456789abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ"

// Generate returns a random base62 code of length n via crypto/rand; the byte%62 modulo bias is negligible.
func Generate(n int) (string, error) {
	b := make([]byte, n)
	if _, err := rand.Read(b); err != nil {
		return "", pkgerrors.WithStack(err)
	}
	for i := range b {
		b[i] = alphabet[int(b[i])%len(alphabet)]
	}
	return string(b), nil
}
