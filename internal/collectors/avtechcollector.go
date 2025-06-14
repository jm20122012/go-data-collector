package collectors

import (
	"context"
	"go-data-collector/internal/devices"
	"time"
)

type AvtechCollector struct {
	CollectorGroupConfig
	Devices []*devices.AvtechSensor // List of Avtech devices to collect data from
}

// NewAvtechCollector creates a new AvtechCollector instance
func NewAvtechCollector(config CollectorGroupConfig, devices []*devices.AvtechSensor) *AvtechCollector {
	return &AvtechCollector{
		CollectorGroupConfig: config,
		Devices:              devices,
	}
}

// Start begins the collection process for Avtech devices
func (c *AvtechCollector) Start() error {
	c.Logger.Info("Starting Avtech collector...")

	// Collect data from all devices concurrently
	for _, device := range c.Devices {
		c.Wg.Add(1)
		go c.collectData(c.Ctx, device)
	}

	<-c.Ctx.Done()
	return nil
}

// Stop halts the collection process for Avtech devices
func (c *AvtechCollector) Stop() error {
	c.Logger.Info("Stopping Avtech collector...")
	c.Wg.Wait() // Wait for all goroutines to finish
	c.Logger.Info("Avtech collector stopped successfully")
	return nil
}

func (c *AvtechCollector) collectData(ctx context.Context, device *devices.AvtechSensor) {
	c.Logger.Info("collecting avtech sensor data", "deviceName", device.Name, "deviceIP", device.IP, "pollInterval", device.PollInterval)

	defer c.Wg.Done()

	// Create a ticker to simulate periodic work/checks
	ticker := time.NewTicker(device.PollInterval)
	defer ticker.Stop()

	device.FetchData()

	for {
		select {
		case <-ctx.Done():
			return
		case <-ticker.C:
			device.FetchData()
		}
	}
}
