package app

import (
	"context"
	"fmt"
	"log/slog"
	"os"
	"os/signal"
	"syscall"
	"time"
)

const shutdownTimeout = 10 * time.Second

// Service is implemented by every runnable binary: Run blocks, Shutdown drains.
type Service interface {
	Run() error
	Shutdown(ctx context.Context) error
}

// RunWithGracefulShutdown runs svc until SIGINT/SIGTERM, then calls Shutdown with a deadline.
func RunWithGracefulShutdown(svc Service) {
	ctx, stop := signal.NotifyContext(context.Background(), syscall.SIGINT, syscall.SIGTERM)
	defer stop()

	go func() {
		if err := svc.Run(); err != nil {
			slog.Error("service crashed", slog.String("error", fmt.Sprintf("%+v", err)))
			os.Exit(1)
		}
	}()

	<-ctx.Done()
	stop()

	shutdownCtx, cancel := context.WithTimeout(context.Background(), shutdownTimeout)
	defer cancel()

	slog.Info("shutting down gracefully")
	if err := svc.Shutdown(shutdownCtx); err != nil {
		slog.Error("graceful shutdown failed", slog.String("error", fmt.Sprintf("%+v", err)))
	}
}
