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

	server := server.NewServer(ctx, cancel, wg, logger, appConfig)
	wg.Add(1)
	server.Start()

	<-ctx.Done() // Wait for context done signal
	wg.Wait()    // Wait for any go routines with a waitgroup to finish
	os.Exit(0)   // Exit cleanly
}

func createLogger(level slog.Level) *slog.Logger {
	handler := slog.NewJSONHandler(os.Stdout, &slog.HandlerOptions{
		Level: level,
	})
	logger := slog.New(handler)
	return logger
}
