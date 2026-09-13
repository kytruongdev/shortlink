package link

import "errors"

// uniqueViolationCode is the Postgres SQLSTATE for a unique-constraint violation.
const uniqueViolationCode = "23505"

var (
	// ErrNotFound is returned when no link matches the query.
	ErrNotFound = errors.New("link not found")
	// ErrConflict is returned on a unique-constraint violation.
	ErrConflict = errors.New("link already exists")
)
