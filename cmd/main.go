package main

import (
	"context"
	cfg "go-data-collector/internal/config"
	"go-data-collector/internal/server"
	"log/slog"
	"os"
	"os/signal"
	"sync"
	"syscall"
)

var (
	// Set at build time.
	version = "dev"
)

func main() {

	slog.Info("Data Collector Server Version", "version", version)

	appConfig := cfg.GetConfig()

	slog.Info("setting logger level", "level", appConfig.LogLevel)
	var logger *slog.Logger
	switch appConfig.LogLevel {
	case "debug":
		logger = createLogger(slog.LevelDebug)
	case "info":
		logger = createLogger(slog.LevelInfo)
	case "warn":
		logger = createLogger(slog.LevelWarn)
	case "error":
		logger = createLogger(slog.LevelError)
	default:
		logger = createLogger(slog.LevelInfo)
	}

	ctx, cancel := context.WithCancel(context.Background())

	c := make(chan os.Signal, 1)
	signal.Notify(c, os.Interrupt, syscall.SIGTERM)

	go func(cancel context.CancelFunc) {
		<-c
		logger.Info("Ctrl+C pressed, cancelling context...")
		cancel()
	}(cancel)

	wg := &sync.WaitGroup{}

	// Initialize the server with the context, wait group, and logger
	server := server.NewServer(ctx, wg, logger)
	server.Start()
	wg.Add(1)

	wg.Wait()
	cancel()
	os.Exit(0)
}

func createLogger(level slog.Level) *slog.Logger {
	handler := slog.NewJSONHandler(os.Stdout, &slog.HandlerOptions{
		Level: level,
	})
	logger := slog.New(handler)
	return logger
}
