package main

import (
	"context"
	"errors"
	"fmt"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/yourusername/go-scaffolding/internal/infrastructure/logger"
)

const (
	defaultConfigPath     = "config.yaml"
	shutdownTimeout       = 30 * time.Second
	serverReadTimeout     = 10 * time.Second
	serverWriteTimeout    = 10 * time.Second
	serverIdleTimeout     = 120 * time.Second
	serverShutdownTimeout = 15 * time.Second
)

func main() {
	// Get config path from environment or use default
	configPath := getConfigPath()

	// Initialize application with all dependencies via Wire
	engine, cleanup, err := initializeApp(configPath)
	if err != nil {
		fmt.Fprintf(os.Stderr, "Failed to initialize application: %v\n", err)
		os.Exit(1)
	}
	defer cleanup()

	// Get port from environment or use default
	port := getPort()

	// Create HTTP server with timeouts
	srv := &http.Server{
		Addr:         fmt.Sprintf(":%s", port),
		Handler:      engine,
		ReadTimeout:  serverReadTimeout,
		WriteTimeout: serverWriteTimeout,
		IdleTimeout:  serverIdleTimeout,
	}

	// Start server in a goroutine
	go func() {
		logger := logger.New("info", os.Stdout)
		logger.Info().
			Str("port", port).
			Str("address", srv.Addr).
			Msg("Starting HTTP server")

		if err := srv.ListenAndServe(); err != nil && !errors.Is(err, http.ErrServerClosed) {
			logger.Fatal().Err(err).Msg("HTTP server failed")
		}
	}()

	// Wait for interrupt signal to gracefully shutdown the server
	quit := make(chan os.Signal, 1)
	signal.Notify(quit, syscall.SIGINT, syscall.SIGTERM)
	<-quit

	// Graceful shutdown
	logger := logger.New("info", os.Stdout)
	logger.Info().Msg("Shutting down server...")

	ctx, cancel := context.WithTimeout(context.Background(), serverShutdownTimeout)
	defer cancel()

	if err := srv.Shutdown(ctx); err != nil {
		logger.Fatal().Err(err).Msg("Server forced to shutdown")
	}

	logger.Info().Msg("Server exited")
}

// getConfigPath returns the config file path from environment or default
func getConfigPath() string {
	if path := os.Getenv("CONFIG_PATH"); path != "" {
		return path
	}
	return defaultConfigPath
}

// getPort returns the port from environment or default
func getPort() string {
	if port := os.Getenv("PORT"); port != "" {
		return port
	}
	return "8080"
}
