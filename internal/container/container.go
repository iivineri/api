package container

import (
	"context"
	"iivineri/internal/config"
	"iivineri/internal/database"
	"iivineri/internal/fiber"
	"iivineri/internal/fiber/modules/auth/handler"
	"iivineri/internal/fiber/shared/middleware"
	"iivineri/internal/logger"
	"iivineri/internal/metrics"
	"iivineri/internal/migration"
	"time"
)

type Container struct {
	Config          *config.Config
	Database        database.DatabaseInterface
	Migration       migration.MigrationInterface
	Server          fiber.ServerInterface
	Metrics         metrics.MetricsInterface
	SystemCollector *metrics.SystemCollector
	Logger          *logger.Logger

	// Auth module
	AuthHandler    *handler.AuthHandler
	AuthMiddleware *middleware.AuthMiddleware

	cleanup []func() error
}

func NewContainer(
	config *config.Config,
	database database.DatabaseInterface,
	migration migration.MigrationInterface,
	server fiber.ServerInterface,
	metricsService metrics.MetricsInterface,
	systemCollector *metrics.SystemCollector,
	logger *logger.Logger,
	authHandler *handler.AuthHandler,
	authMiddleware *middleware.AuthMiddleware,
) *Container {
	container := &Container{
		Config:          config,
		Database:        database,
		Migration:       migration,
		Server:          server,
		Metrics:         metricsService,
		SystemCollector: systemCollector,
		Logger:          logger,
		AuthHandler:     authHandler,
		AuthMiddleware:  authMiddleware,
		cleanup:         make([]func() error, 0),
	}

	container.cleanup = append(container.cleanup, func() error {
		database.Close()
		return nil
	})

	container.cleanup = append(container.cleanup, func() error {
		return migration.Close()
	})

	container.cleanup = append(container.cleanup, func() error {
		systemCollector.Stop()
		return nil
	})

	container.cleanup = append(container.cleanup, func() error {
		ctx, cancel := context.WithTimeout(context.Background(), 30*time.Second)
		defer cancel()
		return server.Shutdown(ctx)
	})

	return container
}

func (c *Container) ConnectDatabase(ctx context.Context) error {
	c.Logger.Debug("Connecting to database pool...")
	if err := c.Database.Connect(ctx); err != nil {
		return err
	}

	c.Logger.Debug("Running database health check...")
	return c.Database.HealthCheck(ctx)
}

func (c *Container) HealthCheck(ctx context.Context) error {
	return c.Database.HealthCheck(ctx)
}

func (c *Container) Shutdown() error {
	c.Logger.Debug("Shutting down container...")

	for i := len(c.cleanup) - 1; i >= 0; i-- {
		if err := c.cleanup[i](); err != nil {
			c.Logger.WithError(err).Error("Error during cleanup")
			return err
		}
	}

	c.Logger.Debug("Container shutdown complete")
	return nil
}
