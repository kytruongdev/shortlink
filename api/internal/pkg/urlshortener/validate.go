package urlshortener

import (
	"errors"
	"net"
	"net/url"
	"regexp"
	"strings"
)

// ErrInvalidURL is returned when a URL fails validation.
var ErrInvalidURL = errors.New("invalid url")

var hostnameRegex = regexp.MustCompile(`(?i)^([a-z0-9-]+\.)+[a-z]{2,63}$`)

// Validate reports whether rawURL is an acceptable absolute http(s) URL with a real domain host (no userinfo, IP, or localhost).
func Validate(rawURL string, maxLen int) error {
	if rawURL == "" {
		return ErrInvalidURL
	}
	if len(rawURL) > maxLen {
		return ErrInvalidURL
	}

	u, err := url.ParseRequestURI(rawURL)
	if err != nil {
		return ErrInvalidURL
	}
	if u.Scheme != "http" && u.Scheme != "https" {
		return ErrInvalidURL
	}
	if u.User != nil {
		return ErrInvalidURL
	}

	host := strings.ToLower(u.Hostname())
	if host == "" || host == "localhost" {
		return ErrInvalidURL
	}
	// Reject IP-literal hosts; a shortener should point at a domain name.
	if net.ParseIP(host) != nil {
		return ErrInvalidURL
	}
	if !isValidHostname(host) {
		return ErrInvalidURL
	}

	return nil
}

// isValidHostname reports whether host looks like a domain name: dot-separated labels with an alphabetic TLD.
func isValidHostname(host string) bool {
	if len(host) > 253 {
		return false
	}
	host = strings.TrimSuffix(host, ".")
	if !hostnameRegex.MatchString(host) {
		return false
	}
	for _, label := range strings.Split(host, ".") {
		if len(label) == 0 || len(label) > 63 {
			return false
		}
		if strings.HasPrefix(label, "-") || strings.HasSuffix(label, "-") {
			return false
		}
	}
	return true
}
