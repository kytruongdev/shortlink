package httpserver

import (
	"context"
	"log/slog"
	"net/http"
	"time"
)

const readinessTimeout = 2 * time.Second

// Pinger reports whether a backing dependency (the database) is reachable.
type Pinger interface {
	Ping(ctx context.Context) error
}

// liveness reports process liveness; it always returns 200 while the process serves.
func liveness(w http.ResponseWriter, _ *http.Request) error {
	return WriteJSON(w, http.StatusOK, map[string]string{"status": "ok"})
}

// readiness returns 200 when the DB is reachable, else 503; the error is logged, not exposed.
func readiness(db Pinger, logger *slog.Logger) HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) error {
		ctx, cancel := context.WithTimeout(r.Context(), readinessTimeout)
		defer cancel()

		if err := db.Ping(ctx); err != nil {
			logger.WarnContext(ctx, "readiness check failed", slog.String("error", err.Error()))
			return WriteJSON(w, http.StatusServiceUnavailable, map[string]string{"status": "unavailable"})
		}

		return WriteJSON(w, http.StatusOK, map[string]string{"status": "ready"})
	}
}
