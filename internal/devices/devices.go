package devices

import (
	"context"
	"go-data-collector/internal/db"
	"log/slog"
	"math"
	"time"
)

type IDevice interface {
	FetchData() // Actual implementation that fetches device data
}

type DeviceConfig struct {
	DeviceID     int32         `json:"deviceId"`
	DeviceTypeID int32         `json:"int"`
	IP           string        `json:"ip"`
	Port         int           `json:"port"`
	Name         string        `json:"name"`
	Location     string        `json:"location"`
	PollInterval time.Duration `json:"poll_interval"` // Interval for polling the device
	Timeout      time.Duration `json:"timeout"`       // Timeout for device responses
	Logger       *slog.Logger
	DBStore      *db.Queries
	Ctx          context.Context
}

func roundFloat(val float64, precision int) float64 {
	// If val is 12.345, divide by precision
	p := float64(precision)
	val = val * p // if precision = 2, results in 1234.5

	// Round which goes to nearest whole integer (away from zero)
	rounded := math.Round(val) // results in 1234

	// Divide by p to move decimal place
	final := rounded / p // results in 12.34
	return final
}
