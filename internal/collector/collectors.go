package collector

import (
	"context"
	"go-data-collector/internal/db"
	"log/slog"
	"sync"
)

type ICollectorGroup interface {
	// Start begins the collection process.
	Start() error
	// Stop halts the collection process.
	Stop() error
}

type CollectorGroupConfig struct {
	Ctx     context.Context `json:"ctx"`      // Context for managing cancellation and timeouts
	Wg      *sync.WaitGroup `json:"wg"`       // WaitGroup for managing goroutines
	Logger  *slog.Logger    `json:"logger"`   // Logger for logging messages
	DBStore *db.Queries     `json:"db_store"` // Database connection pool for storing collected data
}
