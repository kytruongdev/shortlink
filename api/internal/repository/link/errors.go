package link

import "errors"

var (
	// ErrNotFound is returned when no link matches the query.
	ErrNotFound = errors.New("link not found")
	// ErrConflict is returned on a unique-constraint violation.
	ErrConflict = errors.New("link already exists")
)
