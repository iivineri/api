package metrics

import (
	"context"
	"runtime"
	"time"

	"iivineri/internal/database"
	"iivineri/internal/logger"
)

type SystemCollector struct {
	metrics  MetricsInterface
	database database.DatabaseInterface
	logger   *logger.Logger
	stopCh   chan struct{}
}

func NewSystemCollector(
	metrics MetricsInterface,
	database database.DatabaseInterface,
	logger *logger.Logger,
) *SystemCollector {
	return &SystemCollector{
		metrics:  metrics,
		database: database,
		logger:   logger,
		stopCh:   make(chan struct{}),
	}
}

// Start begins collecting system metrics at regular intervals
func (sc *SystemCollector) Start(ctx context.Context) {
	ticker := time.NewTicker(30 * time.Second) // Collect every 30 seconds
	defer ticker.Stop()

	sc.logger.Info("Started system metrics collector")

	for {
		select {
		case <-ctx.Done():
			sc.logger.Info("Stopping system metrics collector")
			return
		case <-sc.stopCh:
			sc.logger.Info("System metrics collector stopped")
			return
		case <-ticker.C:
			sc.collectSystemMetrics()
			sc.collectDatabaseMetrics(ctx)
		}
	}
}

// Stop gracefully stops the metrics collector
func (sc *SystemCollector) Stop() {
	close(sc.stopCh)
}

func (sc *SystemCollector) collectSystemMetrics() {
	// Memory metrics
	var m runtime.MemStats
	runtime.ReadMemStats(&m)
	
	sc.metrics.RecordMemoryUsage(m.Alloc)
	sc.metrics.RecordGoroutines(runtime.NumGoroutine())
	
	sc.logger.WithFields(map[string]interface{}{
		"memory_bytes": m.Alloc,
		"goroutines":   runtime.NumGoroutine(),
	}).Debug("Collected system metrics")
}

func (sc *SystemCollector) collectDatabaseMetrics(ctx context.Context) {
	// Test database health and record metrics
	start := time.Now()
	err := sc.database.HealthCheck(ctx)
	duration := time.Since(start)
	
	success := err == nil
	sc.metrics.RecordDBQuery("health_check", duration, success)
	
	if err != nil {
		sc.logger.WithError(err).Warn("Database health check failed")
	} else {
		sc.logger.WithField("duration_ms", duration.Milliseconds()).Debug("Database health check completed")
	}
}