package httpserver

import (
	"fmt"
	"log/slog"
	"net/http"
)

// HandlerFunc is an http handler that returns an error for centralized handling.
type HandlerFunc func(w http.ResponseWriter, r *http.Request) error

// HandlerErr adapts a HandlerFunc to http.HandlerFunc, centralizing error handling.
func HandlerErr(logger *slog.Logger, fn HandlerFunc) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		if err := fn(w, r); err != nil {
			handleError(w, r, err, logger)
		}
	}
}

// handleError logs the error with stack and answers 500.
func handleError(w http.ResponseWriter, r *http.Request, err error, logger *slog.Logger) {
	logger.ErrorContext(r.Context(), "unhandled error", slog.String("error", fmt.Sprintf("%+v", err)))
	_ = WriteJSON(w, http.StatusInternalServerError, map[string]string{"error": "internal server error"})
}
