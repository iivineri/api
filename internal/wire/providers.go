package wire

import (
	"iivineri/internal/config"
	"iivineri/internal/container"
	"iivineri/internal/database"
	"iivineri/internal/fiber"
	"iivineri/internal/fiber/modules/auth"
	"iivineri/internal/logger"
	"iivineri/internal/metrics"
	"iivineri/internal/migration"

	"github.com/google/wire"
)

func ProvideConfig() *config.Config {
	return config.NewConfig()
}

func ProvideLogger(cfg *config.Config) *logger.Logger {
	return logger.NewLogger(cfg.App.LogLevel())
}

func ProvideDatabase(cfg *config.Config, logger *logger.Logger) database.DatabaseInterface {
	return database.NewDatabase(cfg.Database, logger)
}

func ProvideMigration(cfg *config.Config, logger *logger.Logger) migration.MigrationInterface {
	return migration.NewMigration(cfg.Database, logger)
}

func ProvideMetrics() metrics.MetricsInterface {
	return metrics.NewMetrics()
}

func ProvideServer(cfg *config.Config, logger *logger.Logger, metricsService metrics.MetricsInterface) fiber.ServerInterface {
	return fiber.NewServer(cfg.App, logger, metricsService)
}

func ProvideSystemCollector(metricsService metrics.MetricsInterface, database database.DatabaseInterface, logger *logger.Logger) *metrics.SystemCollector {
	return metrics.NewSystemCollector(metricsService, database, logger)
}

var ProviderSet = wire.NewSet(
	ProvideConfig,
	ProvideLogger,
	ProvideDatabase,
	ProvideMigration,
	ProvideMetrics,
	ProvideServer,
	ProvideSystemCollector,
	
	// Auth module
	auth.AuthProviderSet,
	
	container.NewContainer,
)
