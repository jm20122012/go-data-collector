package server

import (
	"context"
	"fmt"
	"go-data-collector/internal/collectors"
	"go-data-collector/internal/config"
	"go-data-collector/internal/db"
	"go-data-collector/internal/devices"
	"log/slog"
	"sync"
	"time"

	"github.com/jackc/pgx/v5/pgtype"
	"github.com/jackc/pgx/v5/pgxpool"
)

type Server struct {
	ctx             context.Context
	cancel          context.CancelFunc
	wg              *sync.WaitGroup
	logger          *slog.Logger
	appConfig       *config.AppConfig
	collectorGroups map[string]collectors.ICollectorGroup // Map of collector groups to manage
	dbStore         *db.Queries
}

// NewServer creates a new server instance
func NewServer(ctx context.Context, cancel context.CancelFunc, wg *sync.WaitGroup, logger *slog.Logger, appConfig *config.AppConfig) *Server {
	return &Server{
		ctx:             ctx,
		cancel:          cancel,
		wg:              wg,
		logger:          logger,
		appConfig:       appConfig,
		collectorGroups: make(map[string]collectors.ICollectorGroup),
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

	if len(s.collectorGroups) == 0 {
		s.logger.Warn("No collectors configured, server will not collect any data")
		return
	}

	for groupName, group := range s.collectorGroups {
		s.logger.Info("Starting collector", "groupName", groupName)
		if err := group.Start(); err != nil {
			s.logger.Error("Failed to start collector", "error", err)
			return
		}
		s.wg.Add(1) // Increment the wait group counter for each collector
		go func(c collectors.ICollectorGroup) {
			defer s.wg.Done() // Decrement the wait group counter when done
			// Simulate some work for the collector
			<-s.ctx.Done() // Wait for context cancellation
		}(group)
	}

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
	// First get all collector groups where the collector group enabled flag is set to true
	collectorGroups, err := s.dbStore.GetEnabledCollectorGroups(s.ctx)
	if err != nil {
		s.logger.Error("error fetching collector groups from database", "error", err)
		return
	}

	if len(collectorGroups) == 0 {
		s.logger.Warn("no enabled collector groups exist")
		return
	}

	// For all enabled collector groups, get all enabled devices within that group
	for _, group := range collectorGroups {
		pgInt := pgtype.Int4{
			Int32: int32(group.ID),
			Valid: true,
		}
		deviceList, err := s.dbStore.GetEnabledDevicesByCollectorGroupID(s.ctx, pgInt)
		if err != nil {
			s.logger.Error("error fetching devices from database", "collectorGroup", group, "error", err)
			continue
		}

		if len(deviceList) == 0 {
			s.logger.Warn("no enabled devices found for collector group", "group", group)
			continue
		}

		switch group.GroupName {
		case "avtechSensors":
			s.logger.Debug("adding devices to enabledDevicesMap", "devices", deviceList)
			s.collectorGroups[group.GroupName] = s.configureAvtechSensors(deviceList, group.PollIntervalSeconds.Int32)
		}
	}
}

func (s *Server) configureAvtechSensors(deviceListRows []db.DeviceList, groupPollInterval int32) *collectors.AvtechCollector {
	var deviceList []*devices.AvtechSensor

	for _, row := range deviceListRows {
		s.logger.Debug("creating device", "row", row)

		var d devices.AvtechSensor
		d.DeviceID = row.ID
		d.DeviceTypeID = row.DeviceTypeID
		d.IP = row.IpAddress.String
		d.Port = 999
		d.Name = row.DeviceName
		d.Location = row.Location.String
		d.Timeout = 1 * time.Second
		d.Logger = s.logger
		d.DBStore = s.dbStore
		d.Ctx = s.ctx
		deviceList = append(deviceList, &d)

		// Check if device has a specific poll interval, otherwise use group default
		if row.PollIntervalSeconds.Valid {
			// Device has its own poll interval
			d.PollInterval = time.Duration(row.PollIntervalSeconds.Int32) * time.Second
			s.logger.Debug("using device-specific poll interval", "device", row.DeviceName, "interval", d.PollInterval)
		} else {
			// Use collector group's default poll interval
			d.PollInterval = time.Duration(groupPollInterval) * time.Second
			s.logger.Debug("using group default poll interval", "device", row.DeviceName, "interval", d.PollInterval)
		}
	}

	conf := collectors.CollectorGroupConfig{
		Ctx:     s.ctx,
		Wg:      s.wg,
		Logger:  s.logger,
		DBStore: s.dbStore,
	}
	c := collectors.NewAvtechCollector(conf, deviceList)
	return c
}
