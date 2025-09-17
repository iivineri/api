package metrics

import (
	"runtime"
	"time"

	"github.com/gofiber/fiber/v2"
	"github.com/prometheus/client_golang/prometheus/promhttp"
	"github.com/valyala/fasthttp/fasthttpadaptor"
)

// MetricsMiddleware creates a Fiber middleware for collecting HTTP metrics
func MetricsMiddleware(metrics MetricsInterface) fiber.Handler {
	return func(c *fiber.Ctx) error {
		start := time.Now()
		
		// Increment in-flight requests (if we had the gauge)
		// For now, we'll track via the request counter
		
		// Process request
		err := c.Next()
		
		// Calculate duration
		duration := time.Since(start)
		
		// Get response status
		statusCode := c.Response().StatusCode()
		
		// Record metrics
		metrics.RecordHTTPRequest(
			c.Method(),
			c.Route().Path, // Use route pattern instead of actual path for better grouping
			statusCode,
			duration,
		)
		
		return err
	}
}

// SystemMetricsMiddleware periodically records system metrics
func SystemMetricsMiddleware(metrics MetricsInterface) fiber.Handler {
	// Record system metrics on each request (could be optimized with a ticker)
	return func(c *fiber.Ctx) error {
		// Record memory usage
		var m runtime.MemStats
		runtime.ReadMemStats(&m)
		metrics.RecordMemoryUsage(m.Alloc)
		
		// Record goroutine count
		metrics.RecordGoroutines(runtime.NumGoroutine())
		
		return c.Next()
	}
}

// PrometheusHandler returns a Fiber handler for the /metrics endpoint
func PrometheusHandler() fiber.Handler {
	prometheusHandler := promhttp.Handler()
	
	return func(c *fiber.Ctx) error {
		fasthttpadaptor.NewFastHTTPHandler(prometheusHandler)(c.Context())
		return nil
	}
}

// HealthMetricsMiddleware adds custom metrics to health checks
func HealthMetricsMiddleware(metrics MetricsInterface) fiber.Handler {
	return func(c *fiber.Ctx) error {
		// Only apply to health endpoint
		if c.Path() == "/health" {
			metrics.IncrementUserAction("health_check")
		}
		
		return c.Next()
	}
}