package main

import (
	"context"
	"errors"
	"log/slog"
	"net/http"
	"os"
	"time"

	"github.com/kytruongdev/shortlink/internal/config"
	ctrlauth "github.com/kytruongdev/shortlink/internal/controller/auth"
	ctrllink "github.com/kytruongdev/shortlink/internal/controller/link"
	"github.com/kytruongdev/shortlink/internal/handler"
	"github.com/kytruongdev/shortlink/internal/handler/rest"
	"github.com/kytruongdev/shortlink/internal/infra/app"
	"github.com/kytruongdev/shortlink/internal/infra/db/pg"
	"github.com/kytruongdev/shortlink/internal/infra/httpserver"
	repolink "github.com/kytruongdev/shortlink/internal/repository/link"
	repotoken "github.com/kytruongdev/shortlink/internal/repository/token"
	repouser "github.com/kytruongdev/shortlink/internal/repository/user"

	_ "github.com/kytruongdev/shortlink/internal/docs" // registers the OpenAPI spec served at /swagger
)

const (
	readHeaderTimeout = 5 * time.Second
	readTimeout       = 10 * time.Second
	writeTimeout      = 15 * time.Second
	idleTimeout       = 60 * time.Second
)

// @title       ShortLink API
// @version     1.0
// @description URL shortening service: encode a long URL to a short code, decode it back.
// @BasePath    /api/v1
func main() {
	// Composition root: build each resource here, inject downward.
	cfg := config.MustLoad()

	slog.SetDefault(slog.New(slog.NewJSONHandler(os.Stdout, &slog.HandlerOptions{Level: cfg.LogLevel, AddSource: true})))

	pool := pg.MustConnect(cfg.DatabaseURL)
	defer pool.Close()

	repo := repolink.New(pool)
	ctrl := ctrllink.New(repo)
	restHandler := rest.New(ctrl, cfg.BaseURL)

	authCtrl := ctrlauth.New(repouser.New(pool), repotoken.New(pool), []byte(cfg.JWTSecret), cfg.AccessTokenTTL, cfg.RefreshTokenTTL)
	authHandler := rest.NewAuth(authCtrl, cfg.CookieSecure, cfg.RefreshTokenTTL)

	rtr := handler.New(restHandler, authHandler, []byte(cfg.JWTSecret))
	router := httpserver.New(pool, rtr.Routes)

	srv := &http.Server{
		Addr:              ":" + cfg.Port,
		Handler:           router,
		ReadHeaderTimeout: readHeaderTimeout,
		ReadTimeout:       readTimeout,
		WriteTimeout:      writeTimeout,
		IdleTimeout:       idleTimeout,
	}

	slog.Info("server listening", slog.String("addr", srv.Addr))
	app.RunWithGracefulShutdown(&httpService{srv: srv})
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
