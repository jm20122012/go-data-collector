package collectors

import (
	"context"
	"encoding/json"
	"fmt"
	"go-data-collector/internal/db"
	"go-data-collector/internal/utils"
	"net/http"
	"strconv"
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

	// Collect data from all devices concurrently
	for _, device := range c.devices {
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

func (c *AvtechCollector) collectData(ctx context.Context, device AvtechSensorConfig) {
	c.Logger.Info("collecting avtech sensor data", "deviceName", device.Name, "deviceIP", device.IP)

	defer c.Wg.Done()

	// Create a ticker to simulate periodic work/checks
	ticker := time.NewTicker(device.PollInterval)
	defer ticker.Stop()

	c.getData(device)

	for {
		select {
		case <-ctx.Done():
			return
		case <-ticker.C:
			c.getData(device)
		}
	}
}

func (c *AvtechCollector) getData(device AvtechSensorConfig) {
	c.Logger.Debug("fetching avtech sensor data", "deviceName", device.Name, "deviceIP", device.IP)
	var response AvtechResponse

	url := fmt.Sprintf("http://%s/getData.json", device.IP)
	resp, err := http.Get(url)
	if err != nil {
		c.Logger.Error("error performing http request for avtech sensor", "deviceName", device.Name, "deviceIP", device.IP, "error", err)
		return
	}

	defer resp.Body.Close()

	err = json.NewDecoder(resp.Body).Decode(&response)
	if err != nil {
		c.Logger.Error("error decoding json for avtech sensor", "deviceName", device.Name, "deviceIP", device.IP, "error", err)
		return
	}
	c.processAvtechResponse(device, response)

}

func (c *AvtechCollector) processAvtechResponse(device AvtechSensorConfig, response AvtechResponse) {
	c.Logger.Debug("Processing avtech sensor response")

	var tempF float64
	var tempC float64
	tempF, err := strconv.ParseFloat(response.Sensor[0].TempF, 32)
	if err != nil {
		c.Logger.Error("error converting avtech temp f value to float", "error", err)
	}
	tempC, err = strconv.ParseFloat(response.Sensor[0].TempC, 32)
	if err != nil {
		c.Logger.Error("error converting avtech temp c value to float", "error", err)
	}

	timestamp := utils.GetUTCTimestamp()
	humidity := 0.0

	params := db.WriteAvtechRecordParams{
		Timestamp:    timestamp,
		TempF:        tempF,
		TempC:        tempC,
		Humidity:     humidity,
		DeviceID:     int32(device.DeviceID),
		DeviceTypeID: int32(device.DeviceTypeID),
	}

	c.Logger.Debug("avtech write params built", "params", params)

	err = c.DBStore.WriteAvtechRecord(c.Ctx, params)
	if err != nil {
		c.Logger.Error("error writing avtech record", "error", err)
	}

}
