package httpserver

import (
	"log/slog"
	"net/http"
	"time"

	"github.com/go-chi/chi/v5"
	"github.com/go-chi/chi/v5/middleware"
)

// New builds the root router: middleware, health probes, then service routes (nil if none).
func New(logger *slog.Logger, readinessDB Pinger, registerRoutes func(chi.Router)) *chi.Mux {
	r := chi.NewRouter()
	r.Use(middleware.RequestID)
	r.Use(requestLogger(logger))
	r.Use(middleware.Recoverer)

	r.Get("/healthz", HandlerErr(logger, liveness))
	r.Get("/readyz", HandlerErr(logger, readiness(readinessDB, logger)))

	if registerRoutes != nil {
		r.Group(registerRoutes)
	}

	return r
}

// healthPaths are probes logged at Debug to avoid flooding (hit every few seconds).
var healthPaths = map[string]struct{}{
	"/healthz": {},
	"/readyz":  {},
}

// requestLogger logs method, path, status and latency per request (health at Debug).
func requestLogger(logger *slog.Logger) func(http.Handler) http.Handler {
	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			start := time.Now()
			ww := middleware.NewWrapResponseWriter(w, r.ProtoMajor)

			next.ServeHTTP(ww, r)

			level := slog.LevelInfo
			if _, ok := healthPaths[r.URL.Path]; ok {
				level = slog.LevelDebug
			}

			logger.LogAttrs(r.Context(), level, "request",
				slog.String("method", r.Method),
				slog.String("path", r.URL.Path),
				slog.Int("status", ww.Status()),
				slog.Duration("latency", time.Since(start)),
				slog.String("request_id", middleware.GetReqID(r.Context())),
			)
		})
	}
}
