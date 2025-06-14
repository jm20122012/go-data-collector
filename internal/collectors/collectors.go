package collectors

import (
	"context"
	"go-data-collector/internal/db"
	"log/slog"
	"sync"
	"time"
)

type Collector interface {
	// Start begins the collection process.
	Start() error
	// Stop halts the collection process.
	Stop() error
}

type CollectorConfig struct {
	Ctx     context.Context `json:"ctx"`      // Context for managing cancellation and timeouts
	Wg      *sync.WaitGroup `json:"wg"`       // WaitGroup for managing goroutines
	Logger  *slog.Logger    `json:"logger"`   // Logger for logging messages
	DBStore *db.Queries     `json:"db_store"` // Database connection pool for storing collected data
}

type DeviceConfig struct {
	DeviceID     int           `json:"deviceId"`
	DeviceTypeID int           `json:"int"`
	IP           string        `json:"ip"`
	Port         int           `json:"port"`
	Name         string        `json:"name"`
	Location     string        `json:"location"`
	PollInterval time.Duration `json:"poll_interval"` // Interval for polling the device
	Timeout      time.Duration `json:"timeout"`       // Timeout for device responses
}

type AvtechSensorConfig struct {
	DeviceConfig
}

type AmbienWxStationConfig struct {
	DeviceConfig
}
