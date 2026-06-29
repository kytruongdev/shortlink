package urlshortener

import (
	"net/url"
	"strings"

	pkgerrors "github.com/pkg/errors"
)

// Normalize canonicalizes rawURL for dedup: lowercase scheme/host, strip default port and empty fragment.
func Normalize(rawURL string) (string, error) {
	u, err := url.Parse(rawURL)
	if err != nil {
		return "", pkgerrors.WithStack(err)
	}

	u.Scheme = strings.ToLower(u.Scheme)

	host := strings.ToLower(u.Hostname())
	port := u.Port()
	if (u.Scheme == "http" && port == "80") || (u.Scheme == "https" && port == "443") {
		port = ""
	}
	if port != "" {
		u.Host = host + ":" + port
	} else {
		u.Host = host
	}

	return u.String(), nil
}
