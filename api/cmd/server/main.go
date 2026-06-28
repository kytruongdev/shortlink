package main

import (
	"context"
	"errors"
	"log/slog"
	"net/http"
	"os"
	"time"

	"github.com/kytruongdev/shortlink/internal/config"
	"github.com/kytruongdev/shortlink/internal/infra/app"
	"github.com/kytruongdev/shortlink/internal/infra/db/pg"
	"github.com/kytruongdev/shortlink/internal/infra/httpserver"
)

const readHeaderTimeout = 5 * time.Second

func main() {
	// Composition root: build each resource here, inject downward.
	cfg := config.MustLoad()

	logger := slog.New(slog.NewJSONHandler(os.Stdout, &slog.HandlerOptions{Level: cfg.LogLevel}))

	pool := pg.MustConnect(cfg.DatabaseURL)
	defer pool.Close()

	// No client routes to register; httpserver mounts the health probes.
	router := httpserver.New(logger, pool, nil)

	srv := &http.Server{
		Addr:              ":" + cfg.Port,
		Handler:           router,
		ReadHeaderTimeout: readHeaderTimeout,
	}

	logger.Info("server listening", slog.String("addr", srv.Addr))
	app.RunWithGracefulShutdown(logger, &httpService{srv: srv})
}

// httpService adapts *http.Server to app.Service.
type httpService struct{ srv *http.Server }

// Run starts the HTTP server, treating a graceful close as a clean exit.
func (s *httpService) Run() error {
	if err := s.srv.ListenAndServe(); err != nil && !errors.Is(err, http.ErrServerClosed) {
		return err
	}
	return nil
}

// Shutdown gracefully stops the HTTP server.
func (s *httpService) Shutdown(ctx context.Context) error {
	return s.srv.Shutdown(ctx)
}
