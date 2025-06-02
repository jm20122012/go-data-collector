package server

import (
	"context"
	"go-data-collector/internal/collectors"
	"log/slog"
	"sync"
)

type Server struct {
	// ctx is the context for the server
	ctx        context.Context
	wg         *sync.WaitGroup
	logger     *slog.Logger
	collectors []collectors.Collector // List of collectors to manage
}

// NewServer creates a new server instance
func NewServer(ctx context.Context, wg *sync.WaitGroup, logger *slog.Logger) *Server {
	return &Server{
		ctx:    ctx,
		wg:     wg,
		logger: logger,
	}
}

// Start starts the server
func (s *Server) Start() {
	s.logger.Info("Starting server...")

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

	// Here you would typically clean up resources, stop collectors, etc.
	s.wg.Wait() // Wait for all goroutines to finish
	s.logger.Info("Server stopped successfully")
}
