package main

import (
	"context"
	"log/slog"
	"os"
	"os/signal"
	"syscall"
	"time"

	"medical-webhook/internal/infrastructure/bootstrap"
)

const shutdownTimeout = 30 * time.Second

func main() {
	// Initialize application
	app, cleanup, err := bootstrap.InitializeApp()
	if err != nil {
		slog.Error("Failed to initialize application", "error", err)
		os.Exit(1)
	}
	defer cleanup()

	// Start server in goroutine
	go func() {
		slog.Info("Server starting", "port", app.Config.Port, "env", app.Config.AppEnv)
		if err := app.Start(); err != nil {
			slog.Error("Failed to start server", "error", err)
			os.Exit(1)
		}
	}()

	// Wait for interrupt signal for graceful shutdown
	quit := make(chan os.Signal, 1)
	signal.Notify(quit, os.Interrupt, syscall.SIGTERM)
	sig := <-quit

	slog.Info("Shutdown signal received", "signal", sig.String())

	// Create a deadline for the shutdown
	ctx, cancel := context.WithTimeout(context.Background(), shutdownTimeout)
	defer cancel()

	// Shutdown HTTP server with timeout
	if err := app.ShutdownWithContext(ctx); err != nil {
		slog.Error("Error during server shutdown", "error", err)
	}

	slog.Info("Server stopped gracefully")
}
