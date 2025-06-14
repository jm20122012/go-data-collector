package server

import (
	"context"
	"fmt"
	"go-data-collector/internal/collectors"
	"go-data-collector/internal/config"
	"go-data-collector/internal/db"
	"log/slog"
	"sync"
	"time"

	"github.com/jackc/pgx/v5/pgxpool"
)

type Server struct {
	ctx        context.Context
	cancel     context.CancelFunc
	wg         *sync.WaitGroup
	logger     *slog.Logger
	appConfig  *config.AppConfig
	collectors []collectors.Collector // List of collectors to manage
	dbStore    *db.Queries
}

// NewServer creates a new server instance
func NewServer(ctx context.Context, cancel context.CancelFunc, wg *sync.WaitGroup, logger *slog.Logger, appConfig *config.AppConfig) *Server {
	return &Server{
		ctx:       ctx,
		cancel:    cancel,
		wg:        wg,
		logger:    logger,
		appConfig: appConfig,
	}
}

func (s *Server) Start() {
	s.logger.Info("Starting server...")

	defer s.wg.Done()

	if err := s.setupDatabase(); err != nil {
		s.logger.Error("error setting up database", "error", err)
		return
	}

	s.configureCollectors()

	if len(s.collectors) == 0 {
		s.logger.Warn("No collectors configured, server will not collect any data")
		return
	}

	for _, collector := range s.collectors {
		s.logger.Info("Starting collector", "type", collector)
		if err := collector.Start(); err != nil {
			s.logger.Error("Failed to start collector", "error", err)
			return
		}
		s.wg.Add(1) // Increment the wait group counter for each collector
		go func(c collectors.Collector) {
			defer s.wg.Done() // Decrement the wait group counter when done
			// Simulate some work for the collector
			<-s.ctx.Done() // Wait for context cancellation
		}(collector)
	}

	// Here you would typically start your server logic, such as initializing collectors,
	// setting up routes, etc. For now, we just log that the server has started.
	s.logger.Info("Server started successfully")
}

// Stop stops the server
func (s *Server) Stop() {
	s.logger.Info("Stopping server...")

	s.wg.Wait() // Wait for all goroutines to finish
	s.logger.Info("Server stopped successfully")
}

func (s *Server) setupDatabase() error {
	// Create DB connection pool
	connStr := fmt.Sprintf(
		"postgres://%s:%s@%s:%s/%s?sslmode=disable",
		s.appConfig.DBUser,
		s.appConfig.DBPass,
		s.appConfig.DBHost,
		s.appConfig.DBPort,
		s.appConfig.DBName,
	)

	config, err := pgxpool.ParseConfig(connStr)
	if err != nil {
		s.logger.Error("error parsing database config", "error", err)
		return err
	}

	config.MaxConns = 30
	config.MinConns = 5
	config.MaxConnLifetime = 30 * time.Minute
	config.MaxConnIdleTime = 15 * time.Minute

	// Create pool
	pool, err := pgxpool.NewWithConfig(s.ctx, config)
	if err != nil {
		s.logger.Error("failed to create database connection pool", "error", err)
		return err
	}

	// Test connection
	if err := pool.Ping(context.Background()); err != nil {
		s.logger.Error("failed to ping database", "error", err)
		pool.Close() // Clean up on error
		return err
	}

	s.logger.Info("Database ping successful")

	// Create SQLc DB store
	s.dbStore = db.New(pool)

	// Handle cleanup when context is cancelled
	s.wg.Add(1)
	go func() {
		<-s.ctx.Done()
		pool.Close()
		s.wg.Done()
		s.logger.Info("Database pool closed")
	}()

	return nil
}

func (s *Server) configureCollectors() {
	// Check if avtech collector is enabled, and configure collector if it is
	if s.appConfig.EnableAvtechCollector {
		s.logger.Info("avtech collector enabled - configuring collector")
		deviceList, err := s.dbStore.GetDeviceListByDeviceTypeName(s.ctx, "avtechSensor")
		if err != nil {
			s.logger.Error("error fetching avtech device list", "error", err)
		}

		avtechDevices := make([]collectors.AvtechSensorConfig, 0)
		for _, device := range deviceList {
			newDevice := collectors.AvtechSensorConfig{}

			// Set fields individually
			newDevice.DeviceID = int(device.ID)
			newDevice.DeviceTypeID = int(device.DeviceTypeID)
			newDevice.IP = device.IpAddress.String
			newDevice.Port = -1
			newDevice.Name = device.DeviceName
			newDevice.Location = device.Location.String
			newDevice.PollInterval = 10 * time.Minute
			newDevice.Timeout = 1 * time.Second

			avtechDevices = append(avtechDevices, newDevice)
		}

		config := collectors.CollectorConfig{
			Ctx:     s.ctx,
			Wg:      s.wg,
			Logger:  s.logger,
			DBStore: s.dbStore,
		}
		newAvtechCollector := collectors.NewAvtechCollector(config, avtechDevices)

		s.collectors = append(s.collectors, newAvtechCollector)
	}
}
