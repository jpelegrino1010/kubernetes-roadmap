// Package main is the composition root for the API server.
// All dependencies are wired here; no package initialises itself via init().
package main

import (
	"context"
	"errors"
	"fmt"
	"log/slog"
	"net"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/example/app-go/internal/config"
	"github.com/example/app-go/internal/health"
	"github.com/example/app-go/internal/user"
	"github.com/example/app-go/internal/version"
)

func main() {
	cfg := config.Load()

	logLevel := slog.LevelInfo
	if cfg.LogLevel == "debug" {
		logLevel = slog.LevelDebug
	}

	log := slog.New(slog.NewJSONHandler(os.Stdout, &slog.HandlerOptions{Level: logLevel}))
	log.Info("startup", "app_env", cfg.AppEnv, "log_level", cfg.LogLevel)

	if err := run(log); err != nil {
		log.Error("startup", "error", err)
		os.Exit(1)
	}
}

// run contains all startup and shutdown logic so main stays minimal.
// Errors are returned to the caller; nothing is logged here except at the boundary.
func run(log *slog.Logger) error {
	// -------------------------------------------------------------------------
	// Dependencies — constructed explicitly, no global singletons.
	// -------------------------------------------------------------------------
	userStore := user.New()

	// -------------------------------------------------------------------------
	// Router
	// -------------------------------------------------------------------------
	mux := http.NewServeMux()
	mux.Handle("GET /health", health.Handler())
	mux.Handle("GET /users", user.Handler(userStore))
	mux.Handle("GET /version", version.Handler())

	// -------------------------------------------------------------------------
	// Server — validate configuration at startup, not lazily.
	// -------------------------------------------------------------------------
	addr := ":8080"
	srv := &http.Server{
		Addr:         addr,
		Handler:      mux,
		ReadTimeout:  5 * time.Second,
		WriteTimeout: 10 * time.Second,
		IdleTimeout:  120 * time.Second,
		BaseContext: func(l net.Listener) context.Context {
			return context.Background()
		},
	}

	// -------------------------------------------------------------------------
	// Graceful shutdown via OS signals.
	// Every goroutine has a known owner (this function) that manages its lifecycle.
	// -------------------------------------------------------------------------
	shutdown := make(chan os.Signal, 1)
	signal.Notify(shutdown, syscall.SIGINT, syscall.SIGTERM)

	serverErr := make(chan error, 1)

	go func() {
		log.Info("startup", "status", "listening", "addr", addr)
		serverErr <- srv.ListenAndServe()
	}()

	select {
	case err := <-serverErr:
		if !errors.Is(err, http.ErrServerClosed) {
			return fmt.Errorf("run: server: %w", err)
		}

	case sig := <-shutdown:
		log.Info("shutdown", "signal", sig.String())

		ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
		defer cancel()

		if err := srv.Shutdown(ctx); err != nil {
			// Force-close if graceful shutdown times out.
			_ = srv.Close()
			return fmt.Errorf("run: graceful shutdown: %w", err)
		}
	}

	return nil
}
