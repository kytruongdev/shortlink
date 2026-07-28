package user

import "errors"

var (
	// ErrNotFound is returned when no user matches the query.
	ErrNotFound = errors.New("user not found")
	// ErrConflict is returned when the email already exists.
	ErrConflict = errors.New("user already exists")
)
