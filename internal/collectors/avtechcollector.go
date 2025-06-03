package collectors

import (
	"context"
	"time"
)

type AvtechCollector struct {
	CollectorConfig
	devices []AvtechSensorConfig // List of Avtech devices to collect data from
}

// NewAvtechCollector creates a new AvtechCollector instance
func NewAvtechCollector(config CollectorConfig, devices []AvtechSensorConfig) *AvtechCollector {
	return &AvtechCollector{
		CollectorConfig: config,
		devices:         devices,
	}
}

// Start begins the collection process for Avtech devices
func (c *AvtechCollector) Start() error {
	c.Logger.Info("Starting Avtech collector...")
	c.Logger.Info("Avtech collector started successfully")

	for {
		select {
		case <-c.Ctx.Done():
			c.Logger.Info("Context cancelled, stopping collection")
			return c.Stop()
		default:
			// Collect data from all devices concurrently
			for _, device := range c.devices {
				c.Wg.Add(1)

				// Create timeout context for this specific device
				ctx, cancel := context.WithTimeout(c.Ctx, time.Duration(device.Timeout)*time.Second)

				// Launch goroutine with both context and cancel function
				go func(dev AvtechSensorConfig, ctx context.Context, cancelFunc context.CancelFunc) {
					defer c.Wg.Done()
					defer cancelFunc() // Always clean up the context

					c.collectData(ctx, dev)
				}(device, ctx, cancel)
			}

			// Wait for all collection goroutines to finish
			c.Wg.Wait()
		}
	}
}

// Stop halts the collection process for Avtech devices
func (c *AvtechCollector) Stop() error {
	c.Logger.Info("Stopping Avtech collector...")
	c.Wg.Wait() // Wait for all goroutines to finish
	c.Logger.Info("Avtech collector stopped successfully")
	return nil
}

func (c *AvtechCollector) collectData(ctx context.Context, device AvtechSensorConfig) {
	c.Logger.Info("Starting data collection for device: %s", device.Name) // Assuming device has a Name field

	// Create a ticker to simulate periodic work/checks
	ticker := time.NewTicker(100 * time.Millisecond)
	defer ticker.Stop()

	// Simulate data collection work that respects context cancellation
	workDuration := 2 * time.Second // Your simulated work duration
	startTime := time.Now()

	for {
		select {
		case <-ctx.Done():
			// Context was cancelled (either timeout or parent cancellation)
			if ctx.Err() == context.DeadlineExceeded {
				c.Logger.Warn("Data collection timed out for device after %v", time.Since(startTime))
			} else {
				c.Logger.Info("Data collection cancelled for device after %v", time.Since(startTime))
			}
			return

		case <-ticker.C:
			// Check if our simulated work is done
			if time.Since(startTime) >= workDuration {
				c.Logger.Info("Data collection completed successfully for device in %v", time.Since(startTime))
				return
			}
			// Continue working...
		}
	}
}
