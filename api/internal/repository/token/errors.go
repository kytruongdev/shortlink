package token

import "errors"

// ErrNotFound is returned when no refresh token matches the query.
var ErrNotFound = errors.New("refresh token not found")
