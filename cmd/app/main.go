package main

import (
	"log"
	"os"
	"os/signal"
	"syscall"

	"medical-webhook/internal/infrastructure/bootstrap"
)

func main() {
	// Initialize application
	app, cleanup, err := bootstrap.InitializeApp()
	if err != nil {
		log.Fatalf("Failed to initialize application: %v", err)
	}
	defer cleanup()

	// Start server in goroutine
	go func() {
		log.Printf("Server starting on port %s", app.Config.Port)
		if err := app.Start(); err != nil {
			log.Fatalf("Failed to start server: %v", err)
		}
	}()

	// Wait for interrupt signal for graceful shutdown
	quit := make(chan os.Signal, 1)
	signal.Notify(quit, os.Interrupt, syscall.SIGTERM)
	<-quit

	log.Println("Shutting down server...")
	if err := app.Shutdown(); err != nil {
		log.Printf("Error during shutdown: %v", err)
	}

	log.Println("Server stopped gracefully")
}
