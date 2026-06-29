package apperror

import "net/http"

// Error is a domain error carrying an HTTP status and a client-facing code.
type Error struct {
	Code    string
	Message string
	status  int
}

// Error implements the error interface.
func (e *Error) Error() string { return e.Message }

// HTTPStatus returns the HTTP status code for the error.
func (e *Error) HTTPStatus() int { return e.status }

// newError builds an Error; unexported so callers must use the semantic constructors below.
func newError(status int, code, message string) *Error {
	return &Error{Code: code, Message: message, status: status}
}

// BadRequest builds a 400 error.
func BadRequest(code, message string) *Error { return newError(http.StatusBadRequest, code, message) }

// NotFound builds a 404 error.
func NotFound(code, message string) *Error { return newError(http.StatusNotFound, code, message) }

// Internal builds a 500 error.
func Internal(code, message string) *Error {
	return newError(http.StatusInternalServerError, code, message)
}
