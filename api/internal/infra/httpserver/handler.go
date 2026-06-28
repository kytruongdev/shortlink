package httpserver

import "net/http"

// HandlerFunc is an http handler that returns an error for centralized handling.
type HandlerFunc func(w http.ResponseWriter, r *http.Request) error

// HandlerErr adapts a HandlerFunc to http.HandlerFunc, centralizing error handling.
func HandlerErr(fn HandlerFunc) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		if err := fn(w, r); err != nil {
			handleError(w, r, err)
		}
	}
}
