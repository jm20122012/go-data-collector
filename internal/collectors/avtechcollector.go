package collectors

import (
	"context"
	cfg "go-data-collector/internal/config"
	"log/slog"
	"sync"
)

type AvtechCollector struct {
	ctx    context.Context
	wg     *sync.WaitGroup
	logger *slog.Logger
	config *cfg.AppConfig
	name   string // Name of the collector
}

// NewAvtechCollector creates a new AvtechCollector instance
func NewAvtechCollector(ctx context.Context, wg *sync.WaitGroup, logger *slog.Logger, config *cfg.AppConfig) *AvtechCollector {
	return &AvtechCollector{
		ctx:    ctx,
		wg:     wg,
		logger: logger,
		config: config,
		name:   "AvtechCollector", // Name of the collector
	}
}

// Start begins the collection process for Avtech devices
func (c *AvtechCollector) Start() error {
	c.logger.Info("Starting Avtech collector...")

	// Here you would typically start the collection logic for Avtech devices.
	// For now, we just log that the collector has started.
	c.logger.Info("Avtech collector started successfully")

	// Simulate some work
	c.wg.Add(1)
	go func() {
		defer c.wg.Done()
		// Simulate collection work
		<-c.ctx.Done() // Wait for context cancellation
	}()

	return nil
}

// Stop halts the collection process for Avtech devices
func (c *AvtechCollector) Stop() error {
	c.logger.Info("Stopping Avtech collector...")

	// Here you would typically clean up resources, stop collection, etc.
	c.wg.Wait() // Wait for all goroutines to finish
	c.logger.Info("Avtech collector stopped successfully")

	return nil
}

// Name returns the name of the collector
func (c *AvtechCollector) Name() string {
	return c.name
}
