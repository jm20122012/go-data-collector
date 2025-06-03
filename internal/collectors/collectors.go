package collectors

import (
	"context"
	"database/sql"
	"log/slog"
	"sync"
)

type Collector interface {
	// Start begins the collection process.
	Start() error
	// Stop halts the collection process.
	Stop() error
}

type CollectorConfig struct {
	Ctx    context.Context `json:"ctx"`     // Context for managing cancellation and timeouts
	Wg     *sync.WaitGroup `json:"wg"`      // WaitGroup for managing goroutines
	Logger *slog.Logger    `json:"logger"`  // Logger for logging messages
	DBPool *sql.DB         `json:"db_pool"` // Database connection pool for storing collected data
}

type DeviceConfig struct {
	IP           string `json:"ip"`
	Port         int    `json:"port"`
	Name         string `json:"name"`
	Location     string `json:"location"`
	PollInterval int    `json:"poll_interval"` // Interval in seconds for polling the device
	Timeout      int    `json:"timeout"`       // Timeout in seconds for device responses
}

type AvtechSensorConfig struct {
	DeviceConfig
	Username string `json:"username"` // Username for device authentication
	Password string `json:"password"` // Password for device authentication
	URL      string `json:"url"`      // URL for accessing the device's API or web interface
}
