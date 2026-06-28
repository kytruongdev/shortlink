package httpserver

import (
	"context"
	"errors"
	"fmt"
	"log/slog"
	"net/http"

	"github.com/kytruongdev/shortlink/internal/pkg/apperror"
)

type errorResponse struct {
	Code    string `json:"code"`
	Message string `json:"message"`
}

// handleError maps an error to its HTTP response: deadline to 504, *apperror.Error to its own status, else 500.
func handleError(w http.ResponseWriter, _ *http.Request, err error) {
	const (
		codeTimeout  = "TIMEOUT"
		codeInternal = "INTERNAL_ERROR"
	)

	if errors.Is(err, context.DeadlineExceeded) {
		_ = WriteJSON(w, http.StatusGatewayTimeout, errorResponse{Code: codeTimeout, Message: "request timeout"})
		return
	}

	var appErr *apperror.Error
	if errors.As(err, &appErr) {
		_ = WriteJSON(w, appErr.HTTPStatus(), errorResponse{Code: appErr.Code, Message: appErr.Message})
		return
	}

	slog.Error("unhandled error", slog.String("error", fmt.Sprintf("%+v", err)))
	_ = WriteJSON(w, http.StatusInternalServerError, errorResponse{Code: codeInternal, Message: "internal server error"})
}
